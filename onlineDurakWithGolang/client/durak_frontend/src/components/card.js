import React, { Component } from 'react';
import Card from 'react-bootstrap/Card'
import {cardDict} from '../cardDict'

class PlayingCard extends Component {

    render() {
        
        var chosen = this.props.chosenCard? (this.props.chosenCard.Suit === this.props.Suit) && (this.props.chosenCard.Rank === this.props.Rank) : false
        return (
            <Card border='secondary' bg={chosen?"info": ""} text={chosen?"white":""}   onClick={()=> this.props.clicked({Suit: this.props.Suit,Rank: this.props.Rank} )}>
            <Card.Header>{cardDict["Suit"][this.props.Suit]}</Card.Header>
            <Card.Body>
              <Card.Title>{cardDict["Rank"][this.props.Rank]}</Card.Title>
            </Card.Body>
            <Card.Footer> {cardDict["Suit"][this.props.Suit]} </Card.Footer>
          </Card>
        );
    }
}

export default PlayingCard;