# Notion-Meeting-Appのアーキテクチャ
- レイヤードアーキテクチャ + DDD
## 各層とその役割
1. client
    - routerやmiddlewareなど
    - ユーザリクエストを処理する
2. usecase
    - ユーザのリクストに対応する機能など
    - domainに基づいた関数を利用して、高度な機能を提供する
3. domain
    - domain=業務知識
    - domainに基づく基本的な関数を提供する
    - domainを中心とした開発にするために、adapterを利用してinfrastructureがdomainに依存するようにしている=DIPを利用
4. infrastructure
    - DBなど
    - DB操作のインスタンスを提供する
## 依存の流れ
client → usecase → domain-(adapter) ← infrastructure