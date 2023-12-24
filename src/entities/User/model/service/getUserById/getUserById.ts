import { URL } from "../../../const/urls.ts";
import { IUser } from "../../types/userTypes.ts";

export const getUserById = async (id: number): Promise<IUser> => {
    return await fetch( `${ URL }/${ id }` ).then( (res) => res.json() ) as IUser;
};
