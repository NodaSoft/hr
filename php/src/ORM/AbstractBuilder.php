<?php

namespace App\ORM;

abstract class AbstractBuilder
{
    public function __construct(protected readonly QueryBuilder $queryBuilder) {}

    abstract protected function preBuild(): string;

    public function build(): string {

        $sql = $this->preBuild();

        $params = $this->queryBuilder->getParams();
        $conn = $this->queryBuilder->em->connection;

        if($conds = $this->queryBuilder->getConds()) {

            $sql .= " WHERE 1";

            foreach($conds as list($cond, $operator, $col)) {

                $cond = strtoupper($cond);
                $key = ":$col";

                if(is_array($params[$col])) {
                    $in = join(', ', array_map([$conn, 'quote'], $params[$col]));
                    unset($params[$col]);
                    $this->queryBuilder->setParams($params);
                    $key = "($in)";
                    $operator = 'IN';
                }

                $sql .= " $cond `$col` $operator $key ";
            }
        }

        if($orderBy = $this->queryBuilder->getOrderBy()) {
            $orders = [];

            foreach($orderBy as $col => $order) {
                $order = $order === OrderBy::ASC ? 'ASC' : 'DESC';
                $orders[] = "$col $order";
            }

            $sql .= " ORDER BY " . join(', ', $orders);
        }

        if(!is_null($limit = $this->queryBuilder->getLimit())) {
            $sql .= " LIMIT $limit";
        }

        if(!is_null($offset = $this->queryBuilder->getOffset())) {
            $sql .= " OFFSET $offset";
        }

        return $sql;
    }
}