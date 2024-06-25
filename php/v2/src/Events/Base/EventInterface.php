<?php

namespace Nodasoft\Testapp\Events\Base;

interface EventInterface
{
    public function listeners(): array;
}