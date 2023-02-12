<?php

namespace App\ORM;

class DeleteBuilder extends AbstractBuilder
{
    function preBuild(): string
    {
        return "DELETE FROM `{$this->queryBuilder->getFrom()}` WHERE 1";
    }
}