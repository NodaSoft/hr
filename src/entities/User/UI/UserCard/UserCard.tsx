import React from 'react';
import { IUser } from "../../model/types/userTypes.ts";

interface IUserCardProps {
    user: IUser | null;
}

const UserCardComponent: React.FC<IUserCardProps> = (props) => {
    const { user } = props;

    if (!user) {
        return;
    }

    return (
        <table>
            <thead>
            <tr>
                <th>Username</th>
                <th>Phone number</th>
            </tr>
            </thead>
            <tbody>
            <tr>
                <td>{ user.name }</td>
                <td>{ user.phone }</td>
            </tr>
            </tbody>
        </table>
    );
};

export const UserCard = React.memo( UserCardComponent );
