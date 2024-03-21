<?php

    namespace NW\WebService\References\Operations\Notification\BusinessLayer;

    use NW\WebService\References\Operations\Notification\Status;

    class DifferenceMessageRequest
    {


        public int $from;
        public int $to;

        public function __construct(array $raw) {


            $differences = (array)$raw['differences'];

            $this->to = (int)$differences['to'];
            $this->from = (int)$differences['from'];
        }


        public function getFromName():string{
            return Status::getName($this->from);
        }

        public function getToName():string{
            return Status::getName($this->to);
        }
    }
