Feature: draw
    As a leader player, I want to be able to draw the pictionary word
    so that other players can guess.

    Scenario: The leaders draws on the canvas
        Given there is a started game with two players
        When the leader draws the picture
        Then the guesser player sees this picture
