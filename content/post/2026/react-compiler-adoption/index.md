+++
draft = true
title = 'React Compilerの導入意外と大変だった'
categories = ['tech']
tags = ['React']
+++

## はじめに

以前、個人開発のWebアプリにReact Compilerを導入したのですが、依存ライブラリがそこまで多くなかったので、簡単に導入できました。
しかし、副業先では意外と大変だったので共有です。

## 大変だったこと

### `watch()` を `useWatch()` に置き換えた

React Compilerを有効化した後、一部の画面で不具合が見つかり、切り戻すことになりました。。。再導入に向けて既存の実装を見直したとき、React Hook Formの `watch()` が、React公式ドキュメントでメモ化と相性の悪いAPIとして挙げられていました。（なんてこった！）

`watch()` は関数の参照が同じでも、フォームの状態によって異なる値を返します。`watch()` の戻り値を使い回すと、選択肢を変えても入力欄が切り替わりません。そこで、画面表示に使う値は `useWatch()` で購読するように変更しました。

変更は次のようになります。

```diff
- const value = watch('value');
+ const value = useWatch({ control, name: 'value' });
```

一部の画面では、この変更と併せて、画面側の状態管理も修正することになりました。動作確認する範囲も広かったため、修正だけでなく、動作確認まで大変でした。。。

### `"use no memo"` で保存後の入力が反映されない問題を回避した

別のアプリでも、React Compilerを有効化した後に不具合が見つかり、また切り戻すことになりました。再有効化に向けてフォームを調べると、最初の保存は成功するものの、そのまま編集を続けると入力がフォームの値に反映されませんでした。「変更を破棄」の直後にも同じことが起きました。

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

React Hook Formには、`reset()` のあともフィールドの参照を保持する `keepFieldsRef` があります。この設定を使っても、保存や破棄を繰り返すと、入力欄の表示とフォームの値がずれました。

また、React Compilerを有効にしたときの挙動は、ビルド後にしか確認できません。現状ではE2Eテストも存在しないため、動作確認のコストが非常に重い状況です。そこで、当面は `"use no memo"` を使うことにしました。

この例では、親から受け取った `register()` を子が呼ぶので、子コンポーネントの先頭に1行追加しました。

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
- [React Hook Form公式ドキュメント「reset」](https://react-hook-form.com/docs/useform/reset)
