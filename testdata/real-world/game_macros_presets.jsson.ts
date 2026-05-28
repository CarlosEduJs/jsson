export const macros = {
  orbital_strike: {
    category: "Strategem",
    name: "Orbital Strike",
    hotkey: "ctrl+2",
    actions: [
      {
        type: "hold",
        delay: 0,
        key: "ctrl"
      },
      {
        key: "right",
        type: "tap",
        delay: 100
      },
      {
        type: "tap",
        delay: 100,
        key: "right"
      },
      {
        type: "tap",
        delay: 100,
        key: "up"
      },
      {
        type: "release",
        delay: 0,
        key: "ctrl"
      }
    ],
    enabled: true,
    game: "Helldivers 2"
  },
  quick_heal: {
    hotkey: "f1",
    actions: [
      {
        type: "tap",
        delay: 0,
        key: "4"
      },
      {
        delay: 200,
        key: "mouse1",
        type: "tap"
      }
    ],
    game: "Escape from Tarkov",
    category: "Quick Actions",
    enabled: true,
    name: "Quick Heal"
  },
  eagle_airstrike: {
    name: "Eagle Airstrike",
    hotkey: "ctrl+1",
    actions: [
      {
        key: "ctrl",
        type: "hold",
        delay: 0
      },
      {
        delay: 100,
        key: "up",
        type: "tap"
      },
      {
        type: "tap",
        delay: 100,
        key: "right"
      },
      {
        delay: 100,
        key: "down",
        type: "tap"
      },
      {
        type: "release",
        delay: 0,
        key: "ctrl"
      }
    ],
    game: "Helldivers 2",
    category: "Strategem",
    enabled: true
  }
} as const;

export const output = {
  validate: true,
  minify: false,
  format: "json"
} as const;

// Generated types
export type Macros = typeof macros;
export type Output = typeof output;

