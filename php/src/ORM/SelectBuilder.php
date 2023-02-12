<?php

namespace App\ORM;

class SelectBuilder extends AbstractBuilder
{
    protected function preBuild(): string
    {
        return "SELECT * FROM `{$this->queryBuilder->getFrom()}` WHERE 1";
    }
}