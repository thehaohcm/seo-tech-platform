# Web Dashboard

Vue.js 3 frontend for the SEO Tech Platform.

## Features

- 🎨 Modern UI with Tailwind CSS
- 📊 Interactive charts and visualizations
- 🔐 JWT authentication
- 📱 Responsive design
- ⚡ Vite for fast development

## Development

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Project Structure

```
web-dashboard/
├── src/
│   ├── components/      - Reusable Vue components
│   ├── views/           - Page components
│   ├── store/           - Pinia state management
│   ├── router/          - Vue Router configuration
│   └── main.js          - Application entry point
├── public/              - Static assets
└── index.html           - HTML template
```

## Environment Variables

Create a `.env` file:

```
VITE_API_URL=http://localhost:8080
```

## Tech Stack

- Vue 3 (Composition API)
- Vue Router
- Pinia (State Management)
- Axios (HTTP Client)
- Tailwind CSS
- Chart.js
- Vite
