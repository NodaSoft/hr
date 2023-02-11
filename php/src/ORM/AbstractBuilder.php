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

        foreach($this->queryBuilder->getConds() as list($cond, $synbol, $col)) {
            $key = ":$col";

            if(is_array($params[$col])) {
                $in = join(', ', array_map([$conn, 'quote'], $params[$col]));
                unset($params[$col]);
                $this->queryBuilder->setParams($params);
                $key = "($in)";
                $synbol = 'IN';
            }

            $sql .= " $cond `$col` $synbol $key ";
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