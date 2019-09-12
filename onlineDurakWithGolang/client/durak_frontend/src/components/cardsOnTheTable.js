import React, { Component } from 'react';
import '../App.css';
import PlayingCard from './card'
import Row from 'react-bootstrap/Row'
import Container from 'react-bootstrap/Container'

class Table extends Component {

    render() {
        var cardsBeaten = []
        var cardsNotBeaten = []
        for (var i = 0; i < this.props.beaten.length; i++) {
            cardsBeaten.push(
                <PlayingCard key={i+40} Suit={this.props.beaten[i].Suit} Rank={this.props.beaten[i].Rank} clicked={function () { }}> </PlayingCard>
            )
        }
        for (i = 0; i < this.props.notBeaten.length; i++) {
            cardsNotBeaten.push(<PlayingCard key={i+20} Suit={this.props.notBeaten[i].Suit} Rank={this.props.notBeaten[i].Rank} clicked={this.props.action}> </PlayingCard>)
        } 
        return (
            <Container style={{ height: '90%'}}>
                <div>
                Not Beaten cards
                </div>
                <Row style={{ height: '50%', borderStyle:'groove' }}>
                    {cardsNotBeaten}
                </Row>
                <div>
                Beaten cards
                </div>
                <Row style={{ height: '50%', borderStyle:'groove'}}>
                    {cardsBeaten}
                </Row>
            </Container>

        );
    }
}

export default Table;