<?php

namespace NodaSoft\DataMapper\Factory;

use NodaSoft\DataMapper\Mapper\Mapper;

class MapperFactory
{
    public function getMapper(string $instance): Mapper
    {
        $mapperNamespace = preg_replace("/Factory$/", "Mapper", __NAMESPACE__);
        $name = $mapperNamespace . $instance . "Mapper";

        if (! class_exists($name)) {
            throw new \Exception("Mapper class $name doesn't exist.");
        }

        $mapper = new $name();

        if (! $mapper instanceof Mapper) {
            throw new \Exception("Should be instance of " . self::class);
        }

        return $mapper;
    }
}
