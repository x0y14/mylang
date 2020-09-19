import enum
import re

class TKN_OPERATIONS:
	def __init__(self, data):
		self.data = data
		super().__init__()

class TKN_NUM:
	def __init__(self, data):
		self.data = data
		super().__init__()

class TokenKind(enum.Enum):
	NUM = enum.auto()
	OPERATIONS = enum.auto()
	WHITESPACE = enum.auto()
	UNKNOWN = enum.auto()


class TokenManager:
	def __init__(self):
		super().__init__()
	
	def is_what_type(self, t):
		# assert(t_type in ['num', 'op', 'whitespace'])

		if t in ['0', '1', '2', '3', '4', '5', '6', '7', '8', '9']:
			return TokenKind.NUM
		elif t in ['+', '-']:
			return TokenKind.OPERATIONS
		elif t in [' ', '\t']:
			return TokenKind.WHITESPACE
		else:
			return TokenKind.UNKNOWN


class Parser:
	def __init__(self, input_):
		self.pos = 0
		self.input = input_
		self.tkn_manager = TokenManager()
		super().__init__()
	
	def is_eof(self):
		return self.pos >= len(self.input)
	
	def get_char(self) -> str:
		return self.input[self.pos]
	
	def consume_char(self) -> str:
		c = self.input[self.pos]
		self.pos += 1
		return c
	
	def consume_while_reg(self, reg) -> str:
		r = ''
		while self.is_eof() == False:
			if re.match(reg, self.get_char()):
				r += self.consume_char()
			else:
				return r

	def consume_whitespace(self):
		self.consume_while_reg('\s')
	
	def parse(self):
		codes = []

		c = self.consume_char()
		assert(self.tkn_manager.is_what_type(c) == TokenKind.OPERATIONS)
		codes.append(TKN_OPERATIONS(c))
		self.consume_whitespace()

		n1 = self.consume_char()
		assert(self.tkn_manager.is_what_type(n1) == TokenKind.NUM)
		codes.append(TKN_NUM(n1))
		self.consume_whitespace()

		n2 = self.consume_char()
		assert(self.tkn_manager.is_what_type(n2) == TokenKind.NUM)
		codes.append(TKN_NUM(n2))
		self.consume_whitespace()
		return codes
	
	def eat_codes(self, codes):
		result = 0
		op = ''
		ns = []
		for c in codes:
			if type(c) == TKN_OPERATIONS:
				op = c.data
			elif type(c) == TKN_NUM:
				ns.append(c.data)

		n1 = int(ns[0])
		n2 = int(ns[1])
		if op == '-':
			n2 *= -1
		result = n1 + n2
		return result



if __name__ == '__main__':
	ps = Parser('- 0 5')
	r = ps.parse()
	i = ps.eat_codes(r)
	print(i)




