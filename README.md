<div align="center">
  <img width=200 alt="logo Rock, Paper, Scissors Elite" src="doc/assets/logo.svg">

  # Rock, Paper, Scissors Elite

  Most competitive Rock, Paper, Scissors game.

</div>


## Dev setup

Start tailwind watcher:
```bash
npx tailwindcss -i ./static/styles/input.css -o ./static/styles/output.css --minify
```


## Flow chart

```mermaid
graph TD
    A[Start username: type: search] --> B[Add to Pool]
    B --> C[Check Player Availability]
    C --> D{Ping valid before starting game?}
    D -- Yes --> E[Create Game]
    D -- No --> F[Await Player Response]
    F --> C
    E --> G[3s Game]
    G --> H{Round Start while Player Score < 3}
    H -- Yes --> I[8s to Select Move]
    H -- No --> J[Display Player Result]
    I --> K{Move Selected?}
    K -- Yes --> L[Save Move]
    K -- No --> M[Time Ended Determine Winner]
    L --> N[Display Waiting on Player]
    N --> O[Update Score]
    M --> P[Save Player Score round += 1]
    P --> H
    J --> Q[Save Player Result]

    %% Notes
    E:::note --> E_note
    G:::note --> G_note
    E_note[From now on if player disconnect he lose]:::note
    G_note[Only allow move]:::note

    class A,B,C,E,G,H,I,K,L,M,N,O,P,Q node;
    class D,F,J decision;

```
