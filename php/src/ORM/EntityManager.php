<?php

namespace App\ORM;

use App\AbstractRepository;
use App\ORM;
use ReflectionClass;
use Exception;
use PDO;

/**
 * @property-read PDO $connection
 */
class EntityManager
{
    /**
     * @var PDO
     */
    private static $_connection;
    private array $persist = [];

    public function __construct()
    {
        if (is_null(static::$_connection)) {
            $dsn = getenv('DB_DSN');
            $user = getenv('DB_USER');
            $password = getenv('DB_PASSWORD');
            static::$_connection = new PDO($dsn, $user, $password);
            static::$_connection->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
        }
    }

    /**
     * @throws Exception
     */
    public function __get(string $name) {
        if($name === 'connection') {
            return static::$_connection;
        }

        throw new Exception(sprintf('%s:%s not exists!', static::class, $name));
    }

    public function queryBuilder(): QueryBuilder {
        return new QueryBuilder($this);
    }

    public function query(string $sql, array $params): \PDOStatement {
        $stmt = $this->connection->prepare($sql);
        $stmt->execute($params);
        return $stmt;
    }

    public function persist(object $user): void
    {
        $reflect = new \ReflectionObject($user);
        $values = [];

        foreach($reflect->getProperties() as $prop) {
            if(!($attrs = $prop->getAttributes(Column::class))) {
                continue;
            }

            /**
             * @var Column $attr
             */
            $attr = $attrs[0]->newInstance();

            if($prop->isInitialized($user)) {
                $values[$prop->getName()] = $prop->getValue($user);
            } else {
                $values[$prop->getName()] = $prop->getDefaultValue();
            }

            if($attr->type == ColumnType::JSON) {
                $values[$prop->getName()] = json_encode($values[$prop->getName()]);
            }
        }

        foreach($reflect->getAttributes(Entity::class) as $attribute) {
            $attrs = $attribute->getArguments();
            $this->persist[$attrs['repository']][] = $values;
        }
    }

    public function flush(): array
    {
        $ids = [];

        $this->transaction(function() use (&$ids) {

            foreach($this->persist as $repository => $values) {

                /**
                 * @var AbstractRepository $repository
                 */
                $repository = new $repository($this);

                foreach($values as $data) {
                    $repository->queryBuilder()->insert($data)->query();
                    $ids[] = $this->lastInsertId();
                }
            }

        });

        $this->persist = [];
        return $ids;
    }

    public function lastInsertId(): string|bool {
        return $this->connection->lastInsertId();
    }

    public function transaction(callable $fn): void {
        try {
            $this->connection->beginTransaction();
            call_user_func($fn);
            $this->connection->commit();
        } catch(\Exception $e) {
            $this->connection->rollBack();
            throw $e;
        }
    }
}