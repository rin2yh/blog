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

React Compilerを有効化した後、一部の画面で不具合が見つかり、切り戻すことになりました。。。

最初に読んだReact Compilerの導入ガイドには、`watch()` の注意点が載っていなかったので、導入時はこの問題を知りませんでした。切り戻したあとに調べて、React公式のESLintルール `incompatible-library` の解説で見つけました。

`watch()` は関数の参照が同じでも、フォームの状態によって異なる値を返します。`watch()` の戻り値を使い回すと、フォームの値を変えても、画面の表示は古い値のままになります。そこで、画面表示に使う値は `useWatch()` で購読するように変更しました。

変更は次のようになります。

```diff
- const value = watch('value');
+ const value = useWatch({ control, name: 'value' });
```

一部の画面では、この変更と併せて、画面側の状態管理も修正することになりました。動作確認する範囲も広かったため、修正だけでなく、動作確認まで大変でした。。。

### `"use no memo"` で保存後の入力が反映されない問題を回避した

別のアプリでも、React Compilerを有効化した後に不具合が見つかり、またしても切り戻すことになりました。再有効化に向けてフォームを動作確認すると、最初の保存は成功するものの、そのまま編集を続けると入力がフォームの値に反映されませんでした。

修正前のフォームを簡略化すると、次のような実装でした。

```tsx
function ExampleForm() {
  const { register, reset, handleSubmit } = useForm({
    defaultValues: { value: "" },
  });

  const onSubmit = handleSubmit(async ({ value }) => {
    // `saveValue`は入力値を保存する関数
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

React Compilerを有効にした状態で保存すると、次のような挙動となりました。

1. 保存後に `reset()` が呼ばれ、入力欄の登録が消える。
2. フォームを再レンダリングするとき、React Compilerが `register("value")` の戻り値を使い回す（`register()` は呼ばれず、入力欄は再登録されない）。
3. 入力すると `onChange` は呼ばれる。入力欄が未登録のため、React Hook Formは新しい値を取り込まない。

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

## おわりに

最終的に、React Compilerの導入自体は完了しましたが、若干スッキリしない結果になりました。React Hook Formを使う箇所では`"use no memo"`で最適化を外すしかないケースが多く、すべてのコンポーネントに適用できたわけではないからです。

今回の導入で大変だった箇所は、コストが見合わないため、当面は`"use no memo"`で乗り切ろうと思っています。v8の対応が安定した暁には、再度対応を見直していきます。

## 参考文献

- [React Compilerの導入ガイド](https://react.dev/learn/react-compiler/installation)
- [React公式ドキュメント「incompatible-library」](https://react.dev/reference/eslint-plugin-react-hooks/lints/incompatible-library)
- [React公式ドキュメント「Incremental Adoption」](https://react.dev/learn/react-compiler/incremental-adoption)
- [React Hook Form公式ドキュメント「reset」](https://react-hook-form.com/docs/useform/reset)
