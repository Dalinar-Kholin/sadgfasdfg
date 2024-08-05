import {IHurtInfoForComp} from "./handleResult/handleResultInterfaces.ts";
import {handleResults} from "./handleResult/handleResults.ts";
import {hurtNames, IItemToSearch} from "../../interfaces.ts";

export function getHurtResult(Ean: string): Promise<IHurtInfoForComp[]> {
    const url = "/api/takePrice?" + new URLSearchParams({ean: Ean});


    let newData: IHurtInfoForComp[] = [];

    return fetch(url, {
        credentials: "include",
        method: "GET",
    }).then(response => {
        if (response.status === 200) {
            return response.json();
        } else {
            throw new Error("Error");
        }
    }).then(data => {
        data.forEach((element: any) => {
            newData.push( handleResults({name: element.hurtName})(element.result));
        });
        return newData;
    }).catch(err => {
        throw new Error(err);
    });
}

export interface IServerMultipleDataResult{
    Ean : string,
    Item : IHurtInfoForComp
    hurtName : hurtNames
}



export async function getMultipleHurtResult(Items: IItemToSearch[]) {
    const map = new Map<string, IServerMultipleDataResult[]>();

    try {
        const response = await fetch("/api/takePrices", {
            credentials: "include",
            method: "POST",
            body: JSON.stringify({Items: Items}),
            headers: {
                "Content-Type": "application/json"
            }
        });

        if (response.status !== 200) {
            throw new Error("Error");
        }

        const data = await response.json();

        data.map((i : any) => {
            i.Result.map((item : any) => {
                const itemArray = map.get(item.Ean);
                const newItem = {
                    Ean: item.Ean,
                    Item: handleResults({name: i.HurtName})(item.Item),
                    hurtName: i.HurtName
                };
                if (itemArray !== undefined) {
                    itemArray.push(newItem);
                } else {
                    map.set(item.Ean, [newItem]);
                }
            });
        });

        return map;
    } catch (err : any) {
        throw new Error(err.message);
    }
}
