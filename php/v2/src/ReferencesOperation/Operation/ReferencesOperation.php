<?php

namespace NodaSoft\ReferencesOperation\Operation;

use NodaSoft\DataMapper\Factory\MapperFactory;
use NodaSoft\ReferencesOperation\Factory\ReferencesOperationFactory;
use NodaSoft\Request\Request;
use NodaSoft\Result\Operation\ReferencesOperation\ReferencesOperationResult;
use NodaSoft\Result\Operation\ReferencesOperation\TsReturnOperationResult;

class ReferencesOperation
{
    /** @var ReferencesOperationFactory $factory */
    private $factory;

    /** @var MapperFactory $mapperFactory */
    private $mapperFactory;

    public function __construct(
        ReferencesOperationFactory $factory,
        Request $request,
        MapperFactory $mapperFactory
    ) {
        $factory->setRequest($request);
        $this->factory = $factory;
        $this->mapperFactory = $mapperFactory;
    }

    /**
     * @throws \Exception
     * @return TsReturnOperationResult
     */
    public function doOperation(): ReferencesOperationResult
    {
        $result = $this->factory->getResult();
        $params = $this->factory->getParams();
        $command = $this->factory->getCommand(
            $result,
            $params,
            $this->mapperFactory
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
