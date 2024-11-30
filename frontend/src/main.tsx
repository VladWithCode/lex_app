import React from 'react'
import { createRoot } from 'react-dom/client'
import './style.css'
import { QueryClientProvider } from '@tanstack/react-query'
import queryClient from './QueryClient'
import Router from './Router'
import { Quit } from '../wailsjs/runtime/runtime'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <QueryClientProvider client={queryClient}>
            <Router />
        </QueryClientProvider>
    </React.StrictMode>
)

// TODO: Find a more idomatic way to write this
document.body.addEventListener('keydown', (e) => {
    if (e.ctrlKey && e.key === 'Q') {
        Quit()
    }
})
