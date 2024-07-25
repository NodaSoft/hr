<?php

function getResellerEmailFrom(): string
{
    return 'contractor@example.com';
}

/**
 * @return array<string>
 */
function getEmailsByPermit($resellerId, $event): array
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}


function __(string $key, ...$translations): string
{
    return '';
}