# domainについて
このアプリで取り扱うdomainについて
## DB
### Workspace
slackのワークスペース情報
- ID
- Name
- CreatedAt
- UpdatedAt
### User
slackのワークスペースに存在するユーザ
- ID
- SlackUserID
- WorkspaceID
- IsAdministrator
- Name
- CreatedAt
- UpdatedAt
### Notion
ユーザのNotionの情報
- ID
- UserID
- Date
- NotionToken
- NotionDatabaseID
- NotionPageContent
## Notion API
詳細は、[このリンク](https://github.com/KazuwoKiwame12/notion-meeting-app/tree/main/golang/domain/model)を参照すること
### Template
NotionAPIに渡す、作成するページ情報
- Parent
- Properties
- Children
### Block
作成するページにおける内容の情報
- Object
- Type
- Paragraph
- Heading1
- Heading2
- Heading3
- BulletedListItem
- NumberedListItem
- ToDo
- Toggle