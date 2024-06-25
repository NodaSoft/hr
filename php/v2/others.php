<?php

namespace NW\WebService\References\Operations\Notification;

function getResellerEmailFrom()
{
    return 'contractor@example.com';
}

/**
 *
 * @todo looks like $event is redundant for this method ignoring it
 *
 * @param $resellerId
 * @param $event
 *
 * @return string[]
 */
function getEmailsByPermit($resellerId, $event = null )
{
    // fakes the method
    return ['someemeil@example.com', 'someemeil2@example.com'];
}