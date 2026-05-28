export const projects = [
  {
    id: 101,
    fullName: "Alpha Launch (Ana)",
    isComplete: true,
    tasks: [
      {
        description: "Code Review",
        isUrgent: true
      },
      {
        description: "Deploy Prod",
        isUrgent: false
      }
    ]
  },
  {
    id: 102,
    fullName: "Beta Refactor (Bruno)",
    isComplete: false,
    tasks: [
      {
        description: "DB Schema Update",
        isUrgent: true
      },
      {
        isUrgent: false,
        description: "API Doc"
      },
      {
        description: "Test Coverage",
        isUrgent: false
      }
    ]
  },
  {
    isComplete: false,
    tasks: [
      {
        description: "Review backlog",
        isUrgent: false
      }
    ],
    id: 103,
    fullName: "Gamma Cleanup (Carlos)"
  }
] as const;

// Generated types
export type Projects = typeof projects;

