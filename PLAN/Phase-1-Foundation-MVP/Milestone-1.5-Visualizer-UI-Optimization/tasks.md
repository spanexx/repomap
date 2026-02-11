# Task: Optimize Visualizer UI

## Component Modularization
- [x] Refactor `Sidebar.tsx` (Task 1.5.1)
  - [x] Extract `FileHeader` component
  - [x] Extract `DependencyList` component
  - [x] Extract `FileStats` component
- [x] Refactor `Chat.tsx` (Task 1.5.2)
  - [x] Extract `useChat` hook for API logic
  - [x] Separate `MessageList` and `MessageInput` components
- [x] Refactor `App.tsx` (Task 1.5.3)
  - [x] Create `useRepoData` hook for fetching/parsing
  - [x] Move state management to Context (if needed) or simplify hooks

## Context Awareness & Planning (Revised)
- [x] Backend: Update `ChatRequest` struct (Task 1.5.4)
  - [x] Add `Context` fields (SelectedNode, ViewMode)
  - [x] Inject context into System Prompt
- [x] Frontend: Pass UI State to API (Task 1.5.4)
  - [x] Capture selected node from `App` state
- [x] Implement Chat History (Task 1.5.5)
  - [x] Support multi-turn conversation context
- [x] Implement Intent Analysis (Task 1.5.x)
  - [x] Backend: Add `AssignIntent` logic
  - [x] Backend: Call `AssignIntent` in `main.go`
  - [x] Backend: Implement Hybrid Analysis (LLM + Static Fallback)
  - [x] Backend: Implement Batch Processing for LLM
- [x] Integrate Chat in Plan View (Task 1.5.x)
  - [x] Frontend: Refactor `Chat.tsx` for embedded mode
  - [x] Frontend: Update `PlanView.tsx` to include embedded Chat

## Feature Enhancements (User Feedback)
- [ ] Improve Responsive Layout (Task 1.5.6)
  - [x] App: Add mobile header and menu toggle
  - [x] Sidebar: Implement mobile overlay/drawer behavior
  - [x] Chat: Optimize for mobile screens
- [ ] Implement Table View (Task 1.5.8)
  - [ ] Create `TableView.tsx` (Intents, Stats)
  - [ ] Add sorting/filtering
  - [ ] Replace "Flow" view in Header
- [ ] Implement Smart Graph (Task 1.5.9)
  - [x] Add `Graph` support for Intent-based coloring
  - [x] Add Toggle in Header (Intent vs Importance)
- [ ] UI Cleanup (Task 1.5.10)
  - [x] Remove "Rank" and "Radial" views as requested
  - [x] Optimize "Cluster" load time (faster stabilization)
  - [x] Swap "Organic" for a spread-out "Simple" view
- [ ] Polish Styling (Task 1.5.7)
  - [ ] Apply premium aesthetics (typography, spacing)
  - [ ] Add loading states and error boundaries

## Verification
- [ ] Verify all UI features (Upload, Graph, Story, Chat) still work
- [ ] Check console for React warnings/errors
