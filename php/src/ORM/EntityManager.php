<?php

namespace App\ORM;

use App\AbstractRepository;
use App\EntityClassMetadata;
use App\EntityMetadata;
use App\Exception\EntityManagerException;
use App\ORM;
use ReflectionClass;
use Exception;
use PDO;
use ReflectionException;

/**
 * @property-read PDO $connection
 */
class EntityManager
{
    private static ?PDO $_connection = null;
    private array $persist = [];

    public function __construct()
    {
        if (is_null(static::$_connection)) {
            $dsn = sprintf("mysql:dbname=%s;host=%s", getenv('DB_HOST'), getenv('DB_NAME'));
            static::$_connection = new PDO($dsn, getenv('DB_USER'), getenv('DB_PASSWORD'));
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

    /**
     * @throws EntityManagerException
     */
    public function persist(object $entity): void
    {
        $reflect = new \ReflectionObject($entity);

        if(!$reflect->getAttributes(Entity::class)) {
            throw new EntityManagerException(sprintf('Class %s is not entity!', $reflect->getName()));
        }

        $this->persist[] = $entity;
    }

    /**
     * @throws Exception
     */
    public function flush(): void
    {
        $this->transaction(function() {
            foreach($this->persist as $entity) {
                $meta = $this->getEntityMetadata($entity);
                $builder = $this->queryBuilder()->from($meta->entityClass, $meta->from);
                $builder->insert($meta->getValues())->query();
                if($idProp = $meta->getIdProperty()) {
                    $idProp->setValue($entity, $this->lastInsertId());
                }
            }
        });

        $this->persist = [];
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

    public function getEntityClassMetadata(string $entityClass): EntityClassMetadata {
        return new EntityClassMetadata(new ReflectionClass($entityClass));
    }

    public function getEntityMetadata(object $entity): EntityMetadata {
        return new EntityMetadata($entity);
    }
}