<?php

namespace NW\WebService\References\Operations\Notification;

use NodaSoft\ReferencesOperation\Factory\TsReturnOperationFactory;
use NodaSoft\ReferencesOperation\Result\ReferencesOperationResult;
use NodaSoft\ReferencesOperation\Result\TsReturnOperationResult;

class TsReturnOperation extends ReferencesOperation
{
    /** @var TsReturnOperationFactory $factory */
    private $factory;

    public function __construct(TsReturnOperationFactory $factory)
    {
        $this->factory = $factory;
    }

    /**
     * @throws \Exception
     * @return TsReturnOperationResult
     */
    public function doOperation(): ReferencesOperationResult
    {
        $result = $this->factory->getResult();
        $params = $this->factory->getParams();
        $command = $this->factory->getCommand($result, $params);

        try {
            $result = $command->execute();
        } catch (\Throwable $th) {
            $somethingWrong = "Something goes wrong while trying execute a command.";
            throw new \Exception($somethingWrong, 500, $th);
        }

        return $result;
    }
}
