import React, { Component } from 'react';
import '../App.css';
import PlayingCard from './card'
import Row from 'react-bootstrap/Row'
import Col from 'react-bootstrap/Col'

class Hand extends Component {

  constructor(props) {
    super(props);
    this.clicked = this.clicked.bind(this);
  }

  clicked(card){
    this.props.action(card)
  }

  render() {
      var cards = []
      for( var i = 0; i < this.props.cards.length;i++){
          cards.push(<Col key={i} md="auto"><PlayingCard Suit={this.props.cards[i].Suit} Rank={this.props.cards[i].Rank} clicked={this.clicked} chosenCard={this.props.chosenCard}> </PlayingCard></Col> )
      } 
    return (     
        <Row>
          {cards}
        </Row> 
    );
  }
}

export default Hand;


