export const projects = [
  {
    fullName: "Alpha Launch (Ana)",
    isComplete: true,
    tasks: [
      {
        isUrgent: true,
        description: "Code Review"
      },
      {
        isUrgent: false,
        description: "Deploy Prod"
      }
    ],
    id: 101
  },
  {
    tasks: [
      {
        description: "DB Schema Update",
        isUrgent: true
      },
      {
        description: "API Doc",
        isUrgent: false
      },
      {
        description: "Test Coverage",
        isUrgent: false
      }
    ],
    id: 102,
    fullName: "Beta Refactor (Bruno)",
    isComplete: false
  },
  {
    isComplete: false,
    tasks: [
      {
        isUrgent: false,
        description: "Review backlog"
      }
    ],
    id: 103,
    fullName: "Gamma Cleanup (Carlos)"
  }
] as const;

// Generated types
export type Projects = typeof projects;

