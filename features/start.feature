Feature: start
    As a player, I want to be able to start the game
    so all members can start the gameplay.

    Scenario: Not enough players to start
        Given there is only one player
        When the user starts the game
        Then the game is not started

    Scenario: Start the game
        Given there are two players in the game
        When the user starts the game
        Then the game is started
        And only leader player receives a random word
