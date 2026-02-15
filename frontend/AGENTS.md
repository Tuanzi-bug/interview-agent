# Frontend Knowledge Base

**Parent:** [../AGENTS.md](../AGENTS.md)

## OVERVIEW
The `frontend/` directory contains the Next.js source code for the interview system's user interface. It focuses on providing a responsive and interactive experience using modern React patterns.

### Tech Stack
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS & Ant Design
- **State Management**: Zustand
- **API Client**: Axios

## STRUCTURE
- `src/app/`: Page definitions and routing (App Router).
- `src/components/`: Reusable UI components.
  - `common/`: Atomic components (Button, Card, etc.).
  - `layout/`: Structural components (Navbar, Footer).
- `src/store/`: Zustand store definitions for global state.
- `src/services/api/`: API client configuration and endpoint services.
- `src/hooks/`: Custom React hooks for shared logic.
- `src/types/`: TypeScript interfaces and types.

## DATA FLOW
To maintain consistency, follow this standard data flow:

1.  **Component**: React components trigger actions.
2.  **Store/Hook**: 
    - Use **Zustand stores** (`src/store/`) for global state (e.g., auth).
    - Use **Custom Hooks** (`src/hooks/`) for complex logic or data fetching.
3.  **API Service**: Define all network requests in `src/services/api/`.

**Example**: `Page` -> `useAuth hook` -> `authStore` -> `api/client.ts`.

## CONVENTIONS
Refer to the [Root AGENTS.md](../AGENTS.md) for global project conventions.

### 1. Styling
- Use **Tailwind CSS** for layout and custom styling. 
- Use **Ant Design** for complex UI elements (modals, forms).

### 2. Client vs. Server Components
- Use `'use client'` ONLY for files that require interactivity (hooks, state).
- Prefer Server Components for data fetching where SEO or initial performance is critical.

### 3. Linting
- Run `npm run lint` before committing.
