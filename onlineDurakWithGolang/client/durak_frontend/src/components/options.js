import React, { Component } from 'react';
import '../App.css';
import Col from 'react-bootstrap/Col'
import Row from 'react-bootstrap/Row'
import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container';

class Options extends Component {

  constructor(props) {
    super(props);
    this.attackOptions = this.attackOptions.bind(this);
    this.defenceOptions = this.defenceOptions.bind(this);
  }

  attackOptions(){
      if (this.props.canBeDoneWithAttack){
        return (<Button variant="success" size="lg" block onClick={this.props.buttonPress}>Finish attack</Button>);
      } else {
        return (<Button variant="success" size="lg" block disabled onClick={this.props.buttonPress}>Finish attack</Button>);
      }
  }

  defenceOptions(){
      return (
        <Button variant="warning" size="lg" block onClick={this.props.buttonPress}> Take cards</Button>
      );
  }

  render() {
    return (
       <Container>
       <Row className="justify-content-md-center">
           <Col sm={4}>
            {this.props.attack? this.attackOptions(): this.defenceOptions()}
           </Col>
        </Row>
       </Container>
    );
  }
}

export default Options;


