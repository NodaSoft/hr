<?php

namespace App\ORM;

use App\Exception\EntityClassMetadataException;
use ReflectionClass;
use Exception;
use PDO;

/**
 * @property-read PDO $connection
 */
class EntityManager
{
    private static ?PDO $_connection = null;

    /**
     * @psalm-var  array<[string, callable]>
     */
    private array $eventHandlers = [];

    /**
     * @var UnitOfWork[]
     */
    private array $units = [];

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

    public function queryBuilder(): QueryBuilder
    {
        return new QueryBuilder($this);
    }

    public function query(string $sql, array $params): \PDOStatement
    {
        $stmt = $this->connection->prepare($sql);
        $stmt->execute($params);
        return $stmt;
    }

    public function find(string $entityName, int $id): object
    {
        if($unit = $this->findUnit($entityName, $id)) {
            return $unit->entity;
        }

        $entity = $this->queryBuilder()->from($entityName)->where('id', $id)->fetchOne();
        $this->addUnit($entity, UnitOfWorkState::MANAGED);
        return $entity;
    }

    public function findOnBy(string $entityName, array $criteria, array $orderBy = []): ?object
    {
        $builder = $this->queryBuilder()->from($entityName);

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        foreach($orderBy as $col => $order) {
            $builder->orderBy($col, $order);
        }

        if($entity = $builder->fetchOne()) {
            if(!($unit = $this->findUnitByEntity($entity))) {
                $unit = $this->addUnit($entity, UnitOfWorkState::MANAGED);
            }

            return $unit->entity;
        }

        return null;
    }

    public function findAll(string $entityName, array $criteria = [], array $orderBy = [], int $limit = null): array
    {
        $builder = $this->queryBuilder()->from($entityName);

        foreach($criteria as $param => $value) {
            $builder->where($param, $value);
        }

        foreach($orderBy as $col => $order) {
            $builder->orderBy($col, $order);
        }

        if($limit) {
            $builder->limit($limit);
        }

        $result = [];

        foreach($builder->fetchAll() as $entity) {
            if(!($unit = $this->findUnitByEntity($entity))) {
                $unit = $this->addUnit($entity, UnitOfWorkState::MANAGED);
            }

            $result[] = $unit->entity;
        }

        return $result;
    }

    /**
     * @throws EntityClassMetadataException
     */
    public function persist(object $entity): void
    {
        $this->addUnit($entity);
    }

    public function onUpdate(callable $fn): void
    {
        $this->eventHandlers[] = ['update', $fn];
    }

    /**
     * @throws Exception
     */
    public function flush(): void
    {
        $this->transaction(function() {
            foreach($this->units as $unitId => $unit) {
                $meta = $unit->entityMetadata;
                $builder = $this->queryBuilder()->from($meta->entityName);

                switch($unit->getState()) {
                    case UnitOfWorkState::NEW:
                        $builder->insert($meta->getValues())->query();
                        if($idProp = $meta->getIdProperty()) {
                            $idProp->setValue($unit->entity, $this->lastInsertId());
                        }
                        $unit->setState(UnitOfWorkState::MANAGED);
                        break;

                    case UnitOfWorkState::MANAGED:
                        if($unit->isModified()) {
                            $idProp = $meta->getIdProperty()->getName();
                            $values = array_filter($meta->getValues(), fn($k) => $k !== $idProp, ARRAY_FILTER_USE_KEY);
                            $builder->update($values, [$idProp => $meta->getIdPropertyValue()])->query();
                            $this->dispatch('update', $unit->entity);
                        }
                        break;

                    case UnitOfWorkState::REMOVED:
                        $builder->delete()->query();
                        unset($this->units[$unitId]);
                        break;
                }
            }
        });
    }

    private function dispatch(string $event, ...$args): void
    {
        foreach($this->eventHandlers as list($ev, $fn)) {
            if($ev === $event) {
                call_user_func_array($fn, $args);
            }
        }
    }

    public function lastInsertId(): string|bool {
        return $this->connection->lastInsertId();
    }

    /**
     * @throws Exception
     */
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

    /**
     * @throws EntityClassMetadataException
     */
    public function getEntityClassMetadata(string $entityClass): EntityClassMetadata
    {
        try {
            return new EntityClassMetadata(new ReflectionClass($entityClass));
        } catch(\ReflectionException $e) {
            throw new EntityClassMetadataException($e->getMessage());
        }
    }

    /**
     * @throws EntityClassMetadataException
     */
    public function getEntityMetadata(object $entity): EntityMetadata
    {
        return new EntityMetadata($entity);
    }

    /**
     * @throws EntityClassMetadataException
     */
    public function hydrationAll($entities): array {
        return array_map(fn($entity) => $this->hydration($entity), $entities);
    }

    /**
     * @throws EntityClassMetadataException
     */
    public function hydration(object $entity): object
    {
        $unit = $this->findUnitByEntity($entity);

        if($unit) {
            return $unit->entity;
        }

        $reflect = new \ReflectionObject($entity);

        foreach($reflect->getProperties() as $prop) {
            if(!($attrs = $prop->getAttributes(Column::class))) {
                continue;
            }

            /**
             * @var Column $attr
             */
            $attr = $attrs[0]->newInstance();

            if($attr->type == ColumnType::JSON and ($json = $prop->getValue($entity))) {
                $prop->setValue($entity, json_decode($json));
            }
        }

        $this->addUnit($entity, UnitOfWorkState::MANAGED);
        return $entity;
    }

    /**
     * @throws EntityClassMetadataException
     */
    private function addUnit(object $entity, ?UnitOfWorkState $state = UnitOfWorkState::NEW): UnitOfWork
    {
        $meta = $this->getEntityMetadata($entity);
        $unit = new UnitOfWork($entity, $this, $state);
        $this->units[$meta->getEntityObjectId()] = $unit;
        return $unit;
    }

    private function findUnit(string $entityName, int $id): ?UnitOfWork
    {
        foreach($this->units as $unit) {
            if($unit->entityMetadata->entityName === $entityName
                && $unit->entityMetadata->getIdPropertyValue() === $id
            ){
                return $unit;
            }
        }

        return null;
    }

    /**
     * @throws EntityClassMetadataException
     */
    private function findUnitByEntity(object $entity): ?UnitOfWork
    {
        $meta = $this->getEntityMetadata($entity);
        return $this->findUnit($meta->entityName, $meta->getIdPropertyValue());
    }
}