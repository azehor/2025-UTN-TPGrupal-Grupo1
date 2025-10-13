export interface Juego {
  id: string;
  empresa: string;
  nombre: string;
  ordenGrafica: number;
  ordenProcesador: number;
  ordenRAM: number;
  imageURL: string;
}

export const juegosHardCodeados: Juego[] = [
  {
    id: "1",
    empresa: "CD Projekt Red",
    nombre: "Cyberpunk 2077",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 8,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/9/9f/Cyberpunk_2077_box_art.jpg",
  },
  {
    id: "2",
    empresa: "Rockstar Games",
    nombre: "Red Dead Redemption 2",
    ordenGrafica: 10,
    ordenProcesador: 9,
    ordenRAM: 8,
    imageURL:
      "https://image.api.playstation.com/cdn/UP1004/CUSA03041_00/Hpl5MtwQgOVF9vJqlfui6SDB5Jl4oBSq.png?w=440",
  },
  {
    id: "3",
    empresa: "FromSoftware",
    nombre: "Elden Ring",
    ordenGrafica: 8,
    ordenProcesador: 8,
    ordenRAM: 7,
    imageURL:
      "https://image.api.playstation.com/vulcan/ap/rnd/202110/2000/YMUoJUYNX0xWk6eTKuZLr5Iw.jpg",
  },
  {
    id: "4",
    empresa: "Valve",
    nombre: "Half-Life: Alyx",
    ordenGrafica: 10,
    ordenProcesador: 9,
    ordenRAM: 9,
    imageURL:
      "https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/546560/capsule_616x353.jpg?t=1673391297",
  },
  {
    id: "5",
    empresa: "Mojang",
    nombre: "Minecraft",
    ordenGrafica: 3,
    ordenProcesador: 4,
    ordenRAM: 4,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/thumb/b/b6/Minecraft_2024_cover_art.png/250px-Minecraft_2024_cover_art.png",
  },
  {
    id: "6",
    empresa: "Ubisoft",
    nombre: "Assassin's Creed Valhalla",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 8,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/pt/e/e9/Assassins_Creed_Valhalla_capa.png",
  },
  {
    id: "7",
    empresa: "Bethesda",
    nombre: "The Elder Scrolls V: Skyrim",
    ordenGrafica: 8,
    ordenProcesador: 7,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/1/15/The_Elder_Scrolls_V_Skyrim_cover.png",
  },
  {
    id: "8",
    empresa: "CD Projekt Red",
    nombre: "The Witcher 3: Wild Hunt",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 8,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/0/0c/Witcher_3_cover_art.jpg",
  },
  {
    id: "9",
    empresa: "Rockstar Games",
    nombre: "Grand Theft Auto V",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 8,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/a/a5/Grand_Theft_Auto_V.png",
  },
  {
    id: "10",
    empresa: "Valve",
    nombre: "Portal 2",
    ordenGrafica: 7,
    ordenProcesador: 6,
    ordenRAM: 6,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/f/f9/Portal2cover.jpg",
  },
  {
    id: "11",
    empresa: "Epic Games",
    nombre: "Fortnite",
    ordenGrafica: 6,
    ordenProcesador: 5,
    ordenRAM: 5,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/a/ae/Fortnite_Save_The_World.jpg",
  },
  {
    id: "12",
    empresa: "Naughty Dog",
    nombre: "The Last of Us Part II",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 8,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/4/4f/TLOU_P2_Box_Art_2.png",
  },
  {
    id: "13",
    empresa: "Capcom",
    nombre: "Resident Evil Village",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/2/2c/Resident_Evil_Village.png",
  },
  {
    id: "14",
    empresa: "Square Enix",
    nombre: "Final Fantasy VII Remake",
    ordenGrafica: 8,
    ordenProcesador: 7,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/c/ce/FFVIIRemake.png",
  },
  {
    id: "15",
    empresa: "CD Projekt Red",
    nombre: "Gwent",
    ordenGrafica: 5,
    ordenProcesador: 4,
    ordenRAM: 4,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/f/fd/Gwent_cover_art.png",
  },
  {
    id: "16",
    empresa: "FromSoftware",
    nombre: "Dark Souls III",
    ordenGrafica: 8,
    ordenProcesador: 7,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/b/bb/Dark_souls_3_cover_art.jpg",
  },
  {
    id: "17",
    empresa: "Valve",
    nombre: "Left 4 Dead 2",
    ordenGrafica: 7,
    ordenProcesador: 6,
    ordenRAM: 6,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/thumb/b/ba/Left4Dead2.jpg/250px-Left4Dead2.jpg",
  },
  {
    id: "18",
    empresa: "Mojang",
    nombre: "Minecraft Dungeons",
    ordenGrafica: 6,
    ordenProcesador: 5,
    ordenRAM: 5,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/7/73/Minecraft_Dungeons_cover_art.jpg",
  },
  {
    id: "19",
    empresa: "Ubisoft",
    nombre: "Far Cry 6",
    ordenGrafica: 9,
    ordenProcesador: 8,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/3/35/Far_cry_6_cover.jpg",
  },
  {
    id: "20",
    empresa: "Bungie",
    nombre: "Destiny 2",
    ordenGrafica: 8,
    ordenProcesador: 7,
    ordenRAM: 7,
    imageURL:
      "https://upload.wikimedia.org/wikipedia/en/0/05/Destiny_2_%28artwork%29.jpg",
  },
];
