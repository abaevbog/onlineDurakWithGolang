import React, { Component } from 'react';
import './App.css';
import Hand from './components/hand'
import Card from 'react-bootstrap/Card'
import Table from './components/cardsOnTheTable'
import Options from './components/options'
import {cardDict} from './cardDict'

class App extends Component {

  constructor(props) {
    super(props);
    this.state = {
      nonBeaten:[ ],
      beaten: [],
      hand:[],
      attacking:null,
      options:[],
      defenceChoices:[],
      cardForDefence: {},
      socket:null,
      kozir: null,
      finished:false
    }
    this.controller = this.controller.bind(this);
    this.attack = this.attack.bind(this);
    this.defend = this.defend.bind(this);
    this.pickCardForDefence = this.pickCardForDefence.bind(this);
    this.finishOrTake = this.finishOrTake.bind(this);
  }


   controller(that, event) {
    try{
      var obj = JSON.parse(event.data)
    } catch(e){
      console.log(e)
      return
    }
    console.log("EVENT: ", obj)
    switch(Object.keys(obj)[0]){
      case "attack":
        that.setState({
          hand: obj["attack"]["deck"],
          attacking: true,
          nonBeaten:obj["attack"]["nonBeaten"],
          beaten:obj["attack"]["beaten"],
          kozir: obj["attack"]["kozir"][0]["Suit"]
        })
        break
      case "defence":
        that.setState({
          hand: obj["defence"]["deck"],
          attacking: false,
          nonBeaten:obj["defence"]["nonBeaten"],
          beaten:obj["defence"]["beaten"],
          cardForDefence:null,
          kozir:obj["defence"]["kozir"][0]["Suit"]
        })
        break
      case "options":
        that.setState({ options: obj["options"],cardForDefence: null})
        break
      case "defenceChoices":
          that.setState({ defenceChoices: obj["defenceChoices"]})
        break
      case "message":
        var mes = obj["message"]
        if (mes === "Broke connection"){
          console.log("Other player broke connection!")
        }else if (mes === "Done!") {
          that.setState({ finished: true, attacking:null})
        } else {
          console.log("RECEIVED MESSAGE ", mes)
        }
        break
      default:
        console.log("UNKONWN MESSAGE",event.data)
    }
  }

  componentDidMount(){
    var webSocket = new WebSocket("ws://localhost:8080/api",);
    console.log("Launched!")
    this.setState({socket:webSocket})
    webSocket.onmessage = (event) => this.controller(this,event)
  }

  getCardId(dictArr,dict){
    for (var ii=0;ii<dictArr.length;ii++){
      if ( (dictArr[ii].Suit === dict.Suit) && (dictArr[ii].Rank === dict.Rank) ){
        return ii
      }
    }
    return -10
  }


  attack(card){
    console.log("attack with",card)
    var cardId = this.getCardId(this.state.hand, card)
    console.log("Card id: ", cardId)
    if (cardId !== -10){
      this.state.socket.send(cardId.toString())
    }
  }

  defend(attackingCard){
    console.log("defend with",attackingCard)
    var cardId = this.getCardId(this.state.nonBeaten,attackingCard)
    if (cardId != -10) {
      this.state.socket.send(cardId.toString())
    }
  }

  pickCardForDefence(chosenCard){
    console.log("picked for defence, ", chosenCard)
    var cardId = this.getCardId(this.state.hand,chosenCard)
    console.log(this.state)
    if (cardId !== -10 && (this.state.cardForDefence == null || this.state.defenceChoices.length == 0)){
      this.setState({
        cardForDefence: chosenCard
      })
      this.state.socket.send(cardId.toString())
  }
  }

  finishOrTake(){
    console.log("Finish or take")
    this.state.socket.send('1000')
  }




  render() { 
    var message
    var kozir = cardDict["Suit"][this.state.kozir]
    if (this.state.finished){
      message = this.state.hand.length === 0 ? <h2>You won!</h2>: <h2>You lost!</h2>
    } else {
      if (this.state.attacking){
        message = <h2>You are the attacker. Kozir suit: {kozir} </h2>
      } else if (this.state.attacking !== null) {
        message = <h2>You are the defender. Kozir suit: {kozir} </h2>
      } else {
        message = <h2>Waiting for another player to join</h2>
      }
    }
    return (
      <div className="App">
        <Card>
          <Card.Header style={{ height: '15%'}}>
              { message } 
              {this.state.attacking === null? null : <Options attack={this.state.attacking} buttonPress={this.finishOrTake} canBeDoneWithAttack={(this.state.nonBeaten.length + this.state.beaten.length !== 0) && (this.state.nonBeaten.length === 0) }>  </Options>}
          </Card.Header>
          <Card.Body>
            <Table beaten = {this.state.beaten} notBeaten={this.state.nonBeaten} action={this.state.cardForDefence!==null?this.defend:function(){}} defenceOptions={this.state.defenceChoices}> </Table>
          </Card.Body>
          <Card.Footer style={{ height: '25%'}}>
            <Hand cards={this.state.hand} action={this.state.attacking?this.attack:this.pickCardForDefence } attacking={this.state.attacking} chosenCard={this.state.cardForDefence}> </Hand>
          </Card.Footer>
        </Card>
      </div> 
    )
  }
}

export default App;


