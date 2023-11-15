<?php

namespace NodaSoft\DataMapper\Entity;

use NodaSoft\DataMapper\EntityInterface\Entity;
use NodaSoft\DataMapper\EntityTrait;
use NodaSoft\GenericDto\Dto\Dto;
use NodaSoft\Messenger;

class Notification implements Entity, Messenger\Content
{
    use EntityTrait\Entity;

    /** @var string */
    private $bodyTemplate;

    /** @var string */
    private $subjectTemplate;

    public function __construct(
        int $id = null,
        string $name = null,
        string $bodyTemplate = null,
        string $subjectTemplate = null
    ) {
        if ($id) $this->setId($id);
        if ($name) $this->setName($name);
        if ($bodyTemplate) $this->setBodyTemplate($bodyTemplate);
        if ($subjectTemplate) $this->setSubjectTemplate($subjectTemplate);
    }

    public function composeMessageSubject(Dto $params): string
    {
        return $this->fillTemplate($this->subjectTemplate, $params);
    }

    public function composeMessageBody(Dto $params): string
    {
        return $this->fillTemplate($this->bodyTemplate, $params);
    }

    public function fillTemplate(string $template, Dto $params): string
    {
        foreach ($params->toArray() as $param => $value) {
            $key = "#$param#";
            if (strpos($template, $key) > 0) {
                $template = str_replace($key, $value, $template);
            }
        }
        return $template;
    }

    public function getBodyTemplate(): string
    {
        return $this->bodyTemplate;
    }

    public function setBodyTemplate(string $bodyTemplate): void
    {
        $this->bodyTemplate = $bodyTemplate;
    }

    public function getSubjectTemplate(): string
    {
        return $this->subjectTemplate;
    }

    public function setSubjectTemplate(string $subjectTemplate): void
    {
        $this->subjectTemplate = $subjectTemplate;
    }
}
