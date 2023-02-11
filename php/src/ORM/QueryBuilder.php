<?php

namespace App\ORM;

use PDO;

class QueryBuilder
{
    private int|null $limit = null;
    private int|null $offset = null;
    private array $conds = [];
    private array $params = [];
    private string $type = SelectBuilder::class;
    private $entityClass;
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
    public function getEntityClass(): string
    {
        return $this->entityClass;
    }


    public function from(string $entityClass, string $from): static
    {
        $this->entityClass = $entityClass;
        $this->from = $from;
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

    public function limit(int $limit): static {
        $this->limit = $limit;
        return $this;
    }

    public function offset(int $offset): static {
        $this->offset = $offset;
        return $this;
    }

    public function where(string $col, mixed $synbolOrValue, mixed $value = null): static {
        return $this->andWhere(...func_get_args());
    }

    public function andWhere(string $col, mixed $synbolOrValue, mixed $value = null): static {
        if(func_num_args() == 2) {
            $value = $synbolOrValue;
            $synbolOrValue = '=';
        }
        $this->params[$col] = $value;
        $this->conds[] = ['and', $synbolOrValue, $col];
        return $this;
    }

    public function orWhere(string $col, mixed $synbolOrValue, mixed $value = null): static {
        if(func_num_args() == 2) {
            $value = $synbolOrValue;
            $synbolOrValue = '=';
        }
        $this->params[$col] = $value;
        $this->conds[] = ['or', $synbolOrValue, $col];
        return $this;
    }

    public function select(): static {
        $this->type = SelectBuilder::class;
        return $this;
    }

    public function update(): static {
        $this->type = UpdateBuilder::class;
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

    public function fetchOne() {
        if($item = $this->query()->fetchObject($this->getEntityClass())) {
            $item = $this->normalizeObject($item);
        }
        return $item;
    }

    public function fetchAll(): array {
        if($items = $this->query()->fetchAll(PDO::FETCH_CLASS, $this->getEntityClass())) {
            $items = array_map(fn($it) => $this->normalizeObject($it), $items);
        }

        return $items;
    }

    public function query(): \PDOStatement {
        return $this->em->query($this->build(), $this->params);
    }

    private function normalizeObject(object $obj): object {
        $reflect = new \ReflectionObject($obj);

        foreach($reflect->getProperties() as $prop) {
            if(!($attrs = $prop->getAttributes(Column::class))) {
                continue;
            }

            /**
             * @var Column $attr
             */
            $attr = $attrs[0]->newInstance();

            if($attr->type == ColumnType::JSON and ($json = $prop->getValue($obj))) {
                $prop->setValue($obj, json_decode($json));
            }
        }

        return $obj;
    }
}