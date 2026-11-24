+++
draft = true
title = 'React Compilerの導入意外と大変だった'
categories = ['tech']
tags = ['React']
+++

## はじめに

以前、個人開発のアプリにReact Compilerを導入したことがあります。依存ライブラリがそこまで多くなかったので、すんなり導入できました。


副業先でReact Compilerの導入を進めたところ、意外と大変だったので共有します。

## 大変だったこと

### `watch()` を `useWatch()` に置き換えた

React Compilerを有効化した後、一部の画面で不具合が見つかり、導入を切り戻しました。再導入に向けて既存の実装を見直したとき、React Hook Formの `watch()` が、React公式ドキュメントでメモ化と相性の悪いAPIとして挙げられていることを確認しました。

`watch()` は関数の参照が同じでも、フォームの状態によって違う値を返します。戻り値を使い回すと、選択肢を変えても、画面には以前の値に対応する入力欄が表示されたままになります。そこで、画面表示に使う値は `useWatch()` で購読するようにしました。

フィールドを1つに絞ると、変更は次のようになります。

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

```text
保存・破棄 → reset() → フィールドの登録が消える
                         ├─ register() を再実行 → 再登録 → 入力が反映される
                         └─ 戻り値を再利用      → 登録されず、入力が反映されない
```

この例では、親が `reset()` を呼び、子がpropsから受け取った `register()` を呼んでいました。そこで、実際に `register()` を呼ぶ子も含め、該当コンポーネントを `"use no memo"` で最適化対象から外しました。

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

## おわりに

React Compilerは導入できました。React Hook Formを使う箇所では`"use no memo"`で最適化を外すしかないケースが多く、すべてのコンポーネントに適用できたわけではありません。

今回の導入で困った箇所については、当面は`"use no memo"`でよいと思っています。v8の対応が安定したら、アップデートして、最適化から外した箇所を見直したいです。

## 参考文献

- [React公式ドキュメント「incompatible-library」](https://react.dev/reference/eslint-plugin-react-hooks/lints/incompatible-library)
- [React公式ドキュメント「Incremental Adoption」](https://react.dev/learn/react-compiler/incremental-adoption)
