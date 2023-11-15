<?php

namespace NodaSoft\Operation\Operation;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\Operation\Factory\OperationFactory;
use NodaSoft\Request\Request;
use NodaSoft\Operation\Result\Result;
use NodaSoft\Operation\Result\NotifyComplaintNewResult;
use NodaSoft\Dependencies\Dependencies;

class Operation
{
    /** @var OperationFactory $factory */
    private $factory;

    /** @var MapperFactory $mapperFactory */
    private $mapperFactory;

    /** @var Dependencies */
    private $dependencies;

    public function __construct(
        Dependencies     $dependencies,
        OperationFactory $factory,
        Request          $request,
        MapperFactory    $mapperFactory
    ) {
        $this->dependencies = $dependencies;
        $factory->setRequest($request);
        $this->factory = $factory;
        $this->mapperFactory = $mapperFactory;
    }

    /**
     * @return NotifyComplaintNewResult
     * @throws \Exception
     */
    public function doOperation(): Result
    {
        $result = $this->factory->getResult();
        $params = $this->factory->getParams();
        $fetchInitialData = $this->factory->getFetchInitialData($this->mapperFactory);

        if (! $params->isValid()) {
            $missingParams = "Required parameter is missing.";
            throw new \Exception($missingParams, 400);
        }

        try {
            $initialData = $fetchInitialData->fetch($params);
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying to fetch initial data.";
            throw new \Exception($somethingWrong, 500, $th);
        }

        $command = $this->factory->getCommand(
            $result,
            $initialData,
            $this->dependencies
        );

        try {
            $result = $command->execute();
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying execute a command.";
            throw new \Exception($somethingWrong, 500, $th);
        }

        return $result;
    }
}
