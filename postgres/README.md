# PostgreSQLのtips
## 関数定義について
- [参考文献1](https://www.postgresql.jp/document/12/html/sql-createfunction.html)
- [参考文献2](https://www.postgresql.jp/document/12/html/plpgsql-trigger.html)
- PL/pgsqlでは、NEWはRECORDデータ型。
```
    この変数は行レベルのトリガでのINSERT/UPDATE操作によって更新された、新しいデータベースの行を保持します。
    文レベルのトリガおよびDELETE操作では、この変数はnullです。
```

## 言語定義
- [参考文献1](https://www.postgresql.jp/document/12/html/plpgsql-overview.html)
- [参考文献2](https://www.postgresql.jp/document/12/html/sql-createlanguage.html)
```
    CREATE FUNCTION trigger_set_timestamp() RETURNS TRIGGER AS $$
    BEGIN
        NEW.updated_at = NOW();
        RETURN NEW;
    END;
    $$ LANGUAGE plpgsql;
```
- LANGUAGEで$$をplpgsqlとして扱えるようにしている

## トリガーの作成
- [参考文献1](https://www.postgresql.jp/document/12/html/sql-createtrigger.html)

## Database URL
- [参考文献1](https://www.postgresql.org/docs/current/libpq-connect.html)