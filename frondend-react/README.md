# Тестовое задание для React разработчика

Перед вами реальный код из нашего проекта. Мы уже так не делаем и хотим понять, что не делаете и вы.

Что нужно сделать:

- Провести рефакторинг в разделе `src/goods-receipts`;
- Убрать все @ts-ignore в разделе `src/goods-receipts`;
- Исправить ошибки в компоненте `src/goods-receipts/components/operation-form`. Должны корректно работать все поля формы, валидация и т.п. Необходимо исправить поведение фокуса на полях формы при использовании табуляции;
- Переписать `src/goods-receipts/store/operations.gr.ts` на `redux-toolkit`;
- Написать в комментарии краткое резюме по коду: назначение кода, сколько времени вы потратили на рефакторинг и что вам хочется сделать с автором кода :)

Рефакторинг `src/core` на ваше усмотрение, но скорее всего вам придется что-то рефакторить и там т.к мы ждем что вы избавите раздел `src/goods-receipts` от всех @ts-ignore.

Мы даем тестовое задание чтобы:

- Уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших коллег;
- Увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
- Снизить число коротких собеседований, когда мы отказываем сразу же.

Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

Как запустить проект?

В терминале выполнить следующие команды:

```bash
nvm use
```

```bash
npm i
```

```bash
npm run webpack:dev
```
