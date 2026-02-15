# Frontend Technology Stack - Quick Reference Card

## Core Stack
| Layer | Technology | Version | Purpose |
|-------|-----------|---------|---------|
| **UI Components** | Ant Design | 5.12.8 | Enterprise UI library with 50+ components |
| **State Mgmt** | Zustand | 4.4.7 | Simple, lightweight state store + persistence |
| **HTTP Client** | Axios | 1.6.5 | HTTP requests with interceptors + token auto-injection |
| **Framework** | Next.js | 14.0.4 | React meta-framework (SSR, SSG, API rewriting) |
| **Styling** | Tailwind CSS | 3.4.1 | Utility-first CSS framework |
| **Language** | TypeScript | 5.3.3 | Type-safe JavaScript (strict mode) |

---

## Architecture Layers

```
┌─────────────────────────────────────────┐
│        Components (UI/Presentation)      │
│  (Ant Design + Tailwind utilities)      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    Custom Hooks (State Abstraction)      │
│  (useAuth, useXxx)                      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    Services (Business Logic)             │
│  (predictionService, etc.)              │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│  Axios Client (HTTP + Interceptors)     │
│  - Auto-inject Bearer token             │
│  - Handle 401 errors                    │
│  - Adaptive timeouts                    │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│        Backend API (REST)                │
└─────────────────────────────────────────┘
```

---

## File Structure

```
frontend/src/
├── app/                    # Next.js pages (App Router)
│   ├── interview/         # Interview flows
│   ├── user/              # User management
│   ├── resume/            # Resume handling
│   └── layout.tsx         # Root layout
│
├── components/            # React components
│   ├── common/            # Reusable components
│   ├── layout/            # Navbar, Footer
│   └── home/              # Home-specific
│
├── services/api/          # API logic
│   ├── client.ts          # Axios instance
│   └── prediction.ts      # API services
│
├── store/                 # Zustand stores
│   └── authStore.ts       # Auth state
│
├── hooks/                 # Custom hooks
│   └── useAuth.ts         # Auth hook
│
├── types/                 # TypeScript definitions
│   └── global.ts          # Global types
│
├── config/                # Configuration
│   └── api.ts             # API endpoints
│
└── utils/                 # Utilities
    └── format.ts          # Helpers
```

---

## Key Technologies

### 1. Ant Design (UI Framework)
```typescript
import { Button, Card, Form, Modal, Table } from 'antd';
import { CheckOutlined, CloseOutlined } from '@ant-design/icons';

// Usage
<Button type="primary" icon={<CheckOutlined />}>Save</Button>
<Card title="Title">Content</Card>
```

### 2. Zustand (State Management)
```typescript
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: (user) => set({ user, isAuthenticated: true }),
      logout: () => set({ user: null, isAuthenticated: false }),
    }),
    { name: 'auth-storage' } // localStorage key
  )
);

// Usage in components
const { user, login } = useAuthStore();
```

### 3. Axios (HTTP Client)
```typescript
import axios from 'axios';

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || '/api',
  timeout: 10000,
});

// Request interceptor - Auto token injection
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - Error handling
apiClient.interceptors.response.use(
  (response) => {
    if (response.data.code === 200) {
      return response.data.data;
    }
    return Promise.reject(response.data);
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
    }
    return Promise.reject(error);
  }
);
```

### 4. Next.js (Framework)
```typescript
// Server Component (default)
export const metadata: Metadata = {
  title: 'Interview Platform',
};

export default function RootLayout({ children }) {
  return <html><body>{children}</body></html>;
}

// Client Component (interactive)
'use client';
import { useState } from 'react';

export default function Component() {
  const [count, setCount] = useState(0);
  return <button onClick={() => setCount(c => c + 1)}>{count}</button>;
}
```

### 5. Tailwind CSS (Styling)
```tsx
// Utility-first approach
<main className="min-h-screen py-8 px-4 bg-gray-50">
  <div className="max-w-7xl mx-auto">
    <h1 className="text-3xl font-bold text-primary mb-6">Title</h1>
    <p className="text-gray-600 leading-relaxed">Content</p>
  </div>
</main>
```

---

## Development Commands

```bash
# Install dependencies
npm install

# Start development server (port 3000)
npm run dev

# Production build
npm run build

# Run production server
npm start

# Lint code (ESLint)
npm run lint

# Format code (Prettier)
npx prettier --write .
```

---

## Environment Variables

```env
# .env.local (development)
NEXT_PUBLIC_API_BASE_URL=http://localhost:8000/api
INTERNAL_API_URL=http://localhost:8888
```

**Note**: Variables starting with `NEXT_PUBLIC_` are sent to browser. Others are server-only.

---

## API Configuration

```typescript
// src/config/api.ts
export const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL || '/api';

export const INTERVIEW_API = {
  START_STREAM: `${API_BASE_URL}/mianshi/stream/start`,
  END_INTERVIEW: `${API_BASE_URL}/mianshi/interview/end`,
  SUBMIT_ANSWER: `${API_BASE_URL}/mianshi/answer/submit`,
};

export const USER_API = {
  LOGIN: `${API_BASE_URL}/user/login`,
  REGISTER: `${API_BASE_URL}/user/register`,
  GET_PROFILE: `${API_BASE_URL}/user/profile`,
};
```

---

## Authentication Flow

1. **Login**: 
   - Submit credentials → Backend returns JWT token
   - Store token in localStorage
   - Zustand authStore updates (user + isAuthenticated)

2. **API Request**:
   - Axios interceptor detects localStorage token
   - Automatically adds `Authorization: Bearer {token}` header
   - Sends request to API

3. **Response**:
   - If 401: Clear token from localStorage, redirect to login
   - If 200: Return data to component
   - Otherwise: Reject with error

---

## Security Features

✅ **Implemented**
- JWT token injection on every request
- Double header strategy (Authorization + X-Auth-Token)
- 401 response cleanup
- Auth-free route detection

❌ **Not Implemented (Recommended)**
- HTTPOnly cookies (instead of localStorage)
- CSRF token handling
- Content Security Policy headers
- Request rate limiting

---

## Performance Optimizations

### Next.js Built-in
- Code splitting per route
- Server components (less client JS)
- Image lazy-loading
- Font optimization
- CSS purging (Tailwind)

### Application Level
- Token caching (localStorage)
- State persistence (Zustand)
- Adaptive timeouts (3min for evaluations)
- Selective query params

---

## Common Patterns

### Creating a Service
```typescript
// src/services/api/myService.ts
import apiClient from './client';

export const myService = {
  getList: async (page: number) => {
    return apiClient.get('/my-endpoint', { params: { page } });
  },
  
  getDetail: async (id: number) => {
    return apiClient.get(`/my-endpoint/${id}`);
  },
  
  create: async (data: any) => {
    return apiClient.post('/my-endpoint', data);
  },
};
```

### Using in Component
```typescript
'use client';
import { useEffect, useState } from 'react';
import { myService } from '@/services/api/myService';

export default function Component() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true);
    myService.getList(1)
      .then(setData)
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      {loading ? 'Loading...' : <pre>{JSON.stringify(data, null, 2)}</pre>}
    </div>
  );
}
```

### Custom Hook
```typescript
// src/hooks/useMyData.ts
import { useEffect, useState } from 'react';
import { myService } from '@/services/api/myService';

export const useMyData = (page: number = 1) => {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    myService.getList(page)
      .then(setData)
      .catch(setError)
      .finally(() => setLoading(false));
  }, [page]);

  return { data, loading, error };
};

// Usage
const { data, loading, error } = useMyData(1);
```

---

## Debugging Tips

### Check Token
```typescript
// In browser console
localStorage.getItem('token')
```

### View Zustand Store
```typescript
// In browser console
import { useAuthStore } from '@/store/authStore';
useAuthStore.getState()
```

### API Diagnostics
- Component: `<BackendHealthCheck />` in home page
- Checks backend connectivity, CORS, endpoints
- Shows detailed error messages

---

## Testing (Not Yet Configured)

Recommended setup:
```bash
# Unit testing
npm install --save-dev jest @testing-library/react @testing-library/jest-dom

# E2E testing
npm install --save-dev cypress
```

Example test:
```typescript
// __tests__/Button.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import Button from '@/components/common/Button';

test('renders and clicks button', async () => {
  const user = userEvent.setup();
  render(<Button>Click me</Button>);
  
  const button = screen.getByRole('button', { name: /click me/i });
  await user.click(button);
  
  expect(button).toBeInTheDocument();
});
```

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Port 3000 in use | `npm run dev -- -p 3001` |
| Styles not applying | Restart dev server, check Tailwind config |
| Token not persisting | Check localStorage enabled, Zustand persist config |
| API 404 errors | Check API base URL, backend running, routes registered |
| Type errors | Run `npm run build`, check tsconfig.json |
| Memory leak warning | Check useEffect cleanup, unsubscribe from stores |

---

## Future Roadmap

- [ ] Migrate to React Query for caching + server state
- [ ] Add Jest + React Testing Library for unit tests
- [ ] Add Cypress for E2E tests
- [ ] Implement Error Boundary component
- [ ] Add Sentry for error tracking
- [ ] Implement retry logic for failed requests
- [ ] Add HTTPOnly cookies + CSRF protection
- [ ] Set up Storybook for component documentation

---

## Resources

- [Next.js Docs](https://nextjs.org/docs)
- [React Docs](https://react.dev)
- [Ant Design Docs](https://ant.design)
- [Zustand Docs](https://github.com/pmndrs/zustand)
- [Axios Docs](https://axios-http.com)
- [Tailwind CSS Docs](https://tailwindcss.com/docs)
- [TypeScript Docs](https://www.typescriptlang.org/docs)

