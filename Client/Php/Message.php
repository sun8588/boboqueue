<?php
class Message{
	const ADD =100;
	const GET=101;
	const DEL=102;
	const ADD_EXPIRED=103;
	const DEC_EXPIRED=104;
	private $_fd;
	/**
	 * 链接socket
	 * @param String $host
	 * @param Int $prot
	 */
	public function connect($host,$prot){
		$this->_fd=socket_create(AF_INET, SOCK_STREAM, SOL_TCP);
		if(!socket_connect($this->_fd, $host,$prot))
		{
			throw new \Exception('message connection fail',10);
			socket_close($this>_fd);
		}
	}
	public function add($key,$value,$expired=0){

		$keyLen=strlen($key);
		$valueLen=strlen($value);
		$data=pack("CNNa{$keyLen}Na{$valueLen}N",'1',self::ADD,$keyLen,$key,$valueLen,$value,$expired);
		return $this->_sendData($data);
	}
	public function get($num=1){
		return $this->_sendData(pack("CNN",'1',self::GET,$num));
	}
	public function del($key){
		$keyLen=strlen($key);
		return $this->_sendData(pack("CNNa".$keyLen,'1',self::DEL,$keyLen,$key));
	}
	public function addExpired($key,$expired){
		$keyLen=strlen($key);
		return $this->_sendData(pack("CNNa{$keyLen}N",'1',self::ADD_EXPIRED,$keyLen,$key,$expired));
	}
	public function decExpired($key,$expired){
		$keyLen=strlen($key);
		return $this->_sendData(pack("CNNa{$keyLen}N",'1',self::DEC_EXPIRED,$keyLen,$key,$expired));
	}
	private function _sendData($binData){
		$binData.=pack('C','\0');
		var_dump($binData);
		if(!socket_send($this->_fd,$binData,strlen($binData),MSG_OOB|MSG_EOF)){
			throw new \Exception('message send fail',11);
			socket_close($this>_fd);
		}
	}
}

$ms=new Message();
$ms->connect('127.0.0.1',9800);
// $ms->add('test','88',0);
// $ms->get(101);
// $ms->del('test');
// $ms->addExpired('test', 100);
$ms->decExpired('test', 1002);
