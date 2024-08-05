import {IItemToSearch} from "../../../../interfaces.ts";

export default function handlePcMarket(input : string) : IItemToSearch[]{
    const items : IItemToSearch[] = input.split("\n").map((line) => {
        const kodMatch = line.match(/Kod{([^}]*)}/);
        const iloscMatch = line.match(/Ilosc{([^}]*)}/);
        const nameMatch = line.match(/Nazwa{([^}]*)}/);

        const kod = kodMatch ? kodMatch[1] : null;
        const property = iloscMatch ? iloscMatch[1] : null;
        const name = nameMatch ? nameMatch[1] : "";

        return (kod && property && name) ? { Ean: kod, Amount: +property, Name: name } : { Ean : "0", Amount: +"0", Name: ""}
    })



    return items.filter(i => {
        return i.Ean !== "0" && i.Amount !== 0
    })
}