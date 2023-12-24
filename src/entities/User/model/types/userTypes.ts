interface IAddressGeo {
    lat: string,
    lng: string
}

interface IAddress {
    street: string,
    suite: string,
    city: string,
    zipcode: string,
    geo: IAddressGeo;
}

interface ICompany {
    name: string,
    catchPhrase: string,
    bs: string
}

export interface IUser {
    id: number;
    email: string;
    name: string;
    phone: string;
    username: string;
    website: string;
    company: ICompany;
    address: IAddress;
}
