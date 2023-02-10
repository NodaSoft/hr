<?php

namespace Model;

use Db\PdoConnection;

class Model
{

    public static function getPdoInstance(): \PDO
    {
        return PdoConnection::getPdoInstance();
    }
}