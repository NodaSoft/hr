<?php

namespace NodaSoft\Request;

use NodaSoft\Dependencies\Dependencies;
use NodaSoft\Operation\Factory\OperationFactory;

class HttpRequest implements Request
{
    /** @var array<string, mixed> */
    private $data;

    /** @var string */
    private $uri;

    public function __construct()
    {
        $this->data = $_REQUEST['data'] ?? [];
        $this->uri = ltrim($_SERVER['REQUEST_URI'], '/');
    }

    /**
     * @param string $key
     * @return mixed
     */
    public function get(string $key)
    {
        return $this->data[$key] ?? null;
    }

    public function getOperationFactory(
        Dependencies $dependencies
    ): OperationFactory {
        $factoryName = $this->composeOperationFactoryClassName($this->uri);

        if (! class_exists($factoryName)) {
            throw new \Exception('Wrong address.');
        }

        /** @var OperationFactory $factory */
        $factory = new $factoryName($dependencies);
        $factory->setDependencies($dependencies);

        return $factory;
    }

    public function composeOperationFactoryClassName(string $uri): string
    {
        $replacedUri = preg_replace("/[_\/]/", " ", $uri);

        if (! is_string($replacedUri)) {
            throw new \Exception("Expected a string from preg_replace.");
        }

        $capitalizedWords = ucwords($replacedUri);
        $words = explode(' ', $capitalizedWords);
        $interfaceReflection = new \ReflectionClass(OperationFactory::class);
        $namespace = $interfaceReflection->getNamespaceName();
        $name = implode("", $words) . "Factory";
        return $namespace . "\\" . $name;
    }
}
