+++
draft = true
title = 'React Compilerの導入意外と大変だった'
categories = ['tech']
tags = ['React']
+++

## はじめに

副業先でReact Compilerの導入を進めました。同じリポジトリにある10アプリが対象です。アプリごとに有効化するPRを作りました。

導入で大変だったのは、依存ライブラリとの互換性です。特にReact Hook Formは、今回使っていたAPIに関してはCompilerへの対応が十分に進んでいないと感じました。フォームのコードを書き換えたり、一部を最適化対象から外したりして導入を進めました。

## 大変だったこと

### フォームの値を読む `watch()` を置き換えた

一部のアプリでは、React Hook Formの `watch()` を使ってフォームの値を画面に反映していました。たとえば、選択肢の値に応じて、表示する入力欄を切り替える処理です。

`watch()` は同じ関数を呼んでも、フォームの内部状態に応じて違う値を返します。関数の参照が同じだからといって戻り値を使い回すと、入力を変えても表示が更新されません。React公式の `incompatible-library` ルールでも、このAPIが例に挙がっています。

表示に使う値は `useWatch()` で購読するようにしました。フィールドを1つに絞ると、変更は次のようになります。

```diff
- const value = watch('value');
+ const value = useWatch({ control, name: 'value' });
```

一部の画面では、フォームの変更とあわせて、画面側の状態管理も修正しました。

### `"use no memo"` で保存後の入力が反映されない問題を回避した

有効化したアプリの一部は、不具合が見つかって切り戻しました。再有効化に向けてフォームの動作を調べると、最初の編集と保存は成功しました。同じ画面で入力し直すと、入力欄の表示は変わりました。ただ、フォームが管理する値は更新されませんでした。「変更を破棄」の直後にも同じことが起きました。

修正前のフォームを簡略化すると、次のような構成です。`saveValue` は入力値を保存する関数です。

```tsx
function ExampleForm() {
  const { register, reset, handleSubmit } = useForm({
    defaultValues: { value: "" },
  });

  const onSubmit = handleSubmit(async ({ value }) => {
    await saveValue(value);
    reset({ value });
  });

  return (
    <form onSubmit={onSubmit}>
      <ExampleInput register={register} />
      <button type="submit">保存</button>
      <button type="button" onClick={() => reset()}>
        変更を破棄
      </button>
    </form>
  );
}

function ExampleInput({
  register,
}: {
  register: UseFormRegister<{ value: string }>;
}) {
  return <input aria-label="値" {...register("value")} />;
}
```

該当するフォームは、保存や破棄のあとに `reset()` を呼んでいました。この実装では、`reset()` でフィールドの登録が消え、次のレンダー時に `register()` を呼ぶことで再登録していました。Compilerが `register()` の戻り値を再利用すると、呼び出しが省かれ、再登録も行われません。


設定の編集画面では、親が `reset()` を呼び、子がpropsから受け取った `register()` を呼んでいました。そこで、実際に `register()` を呼ぶ子も含め、該当コンポーネントを `"use no memo"` で最適化対象から外しました。

修正では、子コンポーネントの先頭に1行追加しました。

```diff
 function ExampleInput({
   register,
 }: {
   register: UseFormRegister<{ value: string }>;
 }) {
+  "use no memo";
   return <input aria-label="値" {...register("value")} />;
 }
```

さらに同じ呼び出し方を他のアプリから探し、4アプリの計22ファイルを修正しました。

このフォームの確認は、一度目の保存で終えると不具合を見逃します。同じ画面で値を変えてもう一度保存し、再読み込み後に二度目の値が残るところまで確認します。破棄の操作がある画面なら、破棄してから編集と保存をやり直す経路も確認します。どちらも普段の操作ですが、今回壊れたのは `reset()` を挟んだあとの入力でした。

## おわりに

React Compilerは導入できました。React Hook Formを使う箇所では`"use no memo"`で最適化を外すしかないケースが多く、すべてのコンポーネントに適用できたわけではありません。

フォームをメモ化するために、動いている実装を無理に変えるほどのメリットは、今のところ感じていません。今回の導入で困った箇所については、当面は`"use no memo"`でよいと思っています。v8の対応が安定したら、アップデートして、最適化から外した箇所を見直したいです。

## 参考文献

- [React公式ドキュメント「incompatible-library」](https://react.dev/reference/eslint-plugin-react-hooks/lints/incompatible-library)
- [React公式ドキュメント「Incremental Adoption」](https://react.dev/learn/react-compiler/incremental-adoption)
