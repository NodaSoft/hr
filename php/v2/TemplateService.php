<?php

namespace NW\WebService\References\Operations\Notification;


use NW\WebService\References\Operations\Notification\Dto\NotificationData;

class TemplateService
{
    public function render(string $template, NotificationData $data): string
    {
        // В реальном приложении здесь был бы код для рендеринга шаблона
        return "Rendered template: $template";
    }
}