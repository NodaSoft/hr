<?php

namespace App\ORM;

use PDO;

class QueryBuilder
{
    private int|null $limit = null;
    private int|null $offset = null;
    /**
     * @var OrderBy[]
     */
    private array $orderBy = [];
    private array $conds = [];
    private array $params = [];
    private string $type = SelectBuilder::class;
    private $entityName;
    private $from;

    public function __construct(public readonly EntityManager $em)
    {
    }

    /**
     * @return string
     */
    public function getFrom(): string
    {
        return $this->from;
    }

    /**
     * @return string
     */
    public function getEntityName(): string
    {
        return $this->entityName;
    }


    public function from(string $entityName): static
    {
        $meta = $this->em->getEntityClassMetadata($entityName);
        $this->entityName = $entityName;
        $this->from = $meta->from;
        return $this;
    }

    /**
     * @return int|null
     */
    public function getLimit(): ?int
    {
        return $this->limit;
    }

    /**
     * @return int|null
     */
    public function getOffset(): ?int
    {
        return $this->offset;
    }

    /**
     * @return array
     */
    public function getConds(): array
    {
        return $this->conds;
    }

    /**
     * @return array
     */
    public function getParams(): array
    {
        return $this->params;
    }

    /**
     * @return string
     */
    public function getType(): string
    {
        return $this->type;
    }

    /**
     * @return array
     */
    public function getOrderBy(): array
    {
        return $this->orderBy;
    }

    public function limit(int $limit): static {
        $this->limit = $limit;
        return $this;
    }

    public function offset(int $offset): static {
        $this->offset = $offset;
        return $this;
    }

    public function where(string $col, mixed $operatorOrValue, mixed $value = null): static {
        return $this->andWhere(...func_get_args());
    }

    public function andWhere(string $col, mixed $operatorOrValue, mixed $value = null): static {
        if(func_num_args() == 2) {
            $value = $operatorOrValue;
            $operatorOrValue = '=';
        }
        $this->params[$col] = $value;
        $this->conds[] = ['and', $operatorOrValue, $col];
        return $this;
    }

    public function orWhere(string $col, mixed $operatorOrValue, mixed $value = null): static {
        if(func_num_args() == 2) {
            $value = $operatorOrValue;
            $operatorOrValue = '=';
        }
        $this->params[$col] = $value;
        $this->conds[] = ['or', $operatorOrValue, $col];
        return $this;
    }

    public function orderBy(string $column, OrderBy $order) : static {
        $this->orderBy[$column] = $order;
        return $this;
    }

    public function select(): static {
        $this->type = SelectBuilder::class;
        return $this;
    }

    public function update(array $data, array $where): static {
        $this->type = UpdateBuilder::class;
        $this->setParams($data);
        foreach($where as $col => $value) {
            $this->where($col, $value);
        }
        return $this;
    }

    public function insert(array $data): static {
        $this->type = InsertBuilder::class;
        $this->setParams($data);
        return $this;
    }

    public function delete(): static {
        $this->type = DeleteBuilder::class;
        return $this;
    }

    public function setParams(array $params): static {
        $this->params = $params;
        return $this;
    }

    public function build(): string {
        return (new $this->type($this))->build();
    }

    public function fetchOne(): object {
        if($item = $this->query()->fetchObject($this->getEntityName())) {
            $item = $this->em->hydration($item);
        }

        return $item;
    }

    public function fetchAll(): array {
        if($items = $this->query()->fetchAll(PDO::FETCH_CLASS, $this->getEntityName())) {
            $items = $this->em->hydrationAll($items);
        }

        return $items;
    }

    public function query(): \PDOStatement {
        return $this->em->query($this->build(), $this->params);
    }
}