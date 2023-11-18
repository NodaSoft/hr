<?php

namespace NodaSoft\Operation;

use NodaSoft\Request\Request;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Dependencies\Dependencies;

class Operation
{
    /** @var Request */
    private $request;

    /** @var Dependencies */
    private $dependencies;

    public function __construct(Dependencies $dependencies)
    {
        $this->request = $dependencies->getRequest();
        $this->dependencies = $dependencies;
    }

    /**
     * @return Result
     * @throws \Exception
     */
    public function doOperation(): Result
    {
        try {
            $factory = $this->request->getOperationFactory($this->dependencies);
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying to handle an address.";
            throw new \Exception($somethingWrong, 400, $th);
        }

        $fetchInitialData = $factory->getFetchInitialData();

        try {
            $initialData = $fetchInitialData->fetch($this->request);
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying to fetch initial data.";
            throw new \Exception($somethingWrong, 500, $th);
        }

        $command = $factory->getCommand();

        try {
            $result = $command->execute($initialData);
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying execute a command.";
            throw new \Exception($somethingWrong, 500, $th);
        }

        return $result;
    }
}
