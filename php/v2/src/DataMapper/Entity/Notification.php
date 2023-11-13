<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;
use NodaSoft\ReferencesOperation\Params\ReferencesOperationParams;

class Notification implements Entity
{
    use EntityTrait\Entity;

    /** @var string */
    private $template;

    public function __construct(
        int $id = null,
        string $name = null,
        string $template = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($template) $this->setTemplate($template);
    }

    public function composeMessage(ReferencesOperationParams $params): string
    {
        $message = $this->getTemplate();
        foreach ($params->toArray() as $param => $value) {
            $key = "#$param#";
            if (strpos($message, $key) > 0) {
                $message = str_replace($key, $value, $message);
            }
        }
        return $message;
    }

    public function getTemplate(): string
    {
        return $this->template;
    }

    public function setTemplate(string $template): void
    {
        $this->template = $template;
    }
}
