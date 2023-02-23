<?php

namespace App\ORM;

enum UnitOfWorkState {
    case NEW;
    case MANAGED;
    case REMOVED;
}
