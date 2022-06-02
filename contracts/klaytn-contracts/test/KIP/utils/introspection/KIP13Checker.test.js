require('@openzeppelin/test-helpers');

const KIP13CheckerMock = artifacts.require('KIP13CheckerMock');
const KIP13NotSupported = artifacts.require('KIP13NotSupported');
const KIP13InterfacesSupported = artifacts.require('KIP13InterfacesSupported');

const DUMMY_ID = '0xdeadbeef';
const DUMMY_ID_2 = '0xcafebabe';
const DUMMY_ID_3 = '0xdecafbad';
const DUMMY_UNSUPPORTED_ID = '0xbaddcafe';
const DUMMY_UNSUPPORTED_ID_2 = '0xbaadcafe';
const DUMMY_ACCOUNT = '0x1111111111111111111111111111111111111111';

contract('KIP13Checker', function () {
  beforeEach(async function () {
    this.mock = await KIP13CheckerMock.new();
  });

  context('KIP13 not supported', function () {
    beforeEach(async function () {
      this.target = await KIP13NotSupported.new();
    });

    it('does not support KIP13', async function () {
      const supported = await this.mock.supportsKIP13(this.target.address);
      expect(supported).to.equal(false);
    });

    it('does not support mock interface via supportsInterface', async function () {
      const supported = await this.mock.supportsInterface(
        this.target.address,
        DUMMY_ID,
      );
      expect(supported).to.equal(false);
    });

    it('does not support mock interface via supportsAllInterfaces', async function () {
      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        [DUMMY_ID],
      );
      expect(supported).to.equal(false);
    });
  });

  context('KIP13 supported', function () {
    beforeEach(async function () {
      this.target = await KIP13InterfacesSupported.new([]);
    });

    it('supports KIP13', async function () {
      const supported = await this.mock.supportsKIP13(this.target.address);
      expect(supported).to.equal(true);
    });

    it('does not support mock interface via supportsInterface', async function () {
      const supported = await this.mock.supportsInterface(
        this.target.address,
        DUMMY_ID,
      );
      expect(supported).to.equal(false);
    });

    it('does not support mock interface via supportsAllInterfaces', async function () {
      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        [DUMMY_ID],
      );
      expect(supported).to.equal(false);
    });
  });

  context('KIP13 and single interface supported', function () {
    beforeEach(async function () {
      this.target = await KIP13InterfacesSupported.new([DUMMY_ID]);
    });

    it('supports KIP13', async function () {
      const supported = await this.mock.supportsKIP13(this.target.address);
      expect(supported).to.equal(true);
    });

    it('supports mock interface via supportsInterface', async function () {
      const supported = await this.mock.supportsInterface(
        this.target.address,
        DUMMY_ID,
      );
      expect(supported).to.equal(true);
    });

    it('supports mock interface via supportsAllInterfaces', async function () {
      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        [DUMMY_ID],
      );
      expect(supported).to.equal(true);
    });
  });

  context('KIP13 and many interfaces supported', function () {
    beforeEach(async function () {
      this.supportedInterfaces = [DUMMY_ID, DUMMY_ID_2, DUMMY_ID_3];
      this.target = await KIP13InterfacesSupported.new(
        this.supportedInterfaces,
      );
    });

    it('supports KIP13', async function () {
      const supported = await this.mock.supportsKIP13(this.target.address);
      expect(supported).to.equal(true);
    });

    it('supports each interfaceId via supportsInterface', async function () {
      for (const interfaceId of this.supportedInterfaces) {
        const supported = await this.mock.supportsInterface(
          this.target.address,
          interfaceId,
        );
        expect(supported).to.equal(true);
      }
    });

    it('supports all interfaceIds via supportsAllInterfaces', async function () {
      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        this.supportedInterfaces,
      );
      expect(supported).to.equal(true);
    });

    it('supports none of the interfaces queried via supportsAllInterfaces', async function () {
      const interfaceIdsToTest = [
        DUMMY_UNSUPPORTED_ID,
        DUMMY_UNSUPPORTED_ID_2,
      ];

      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        interfaceIdsToTest,
      );
      expect(supported).to.equal(false);
    });

    it('supports not all of the interfaces queried via supportsAllInterfaces', async function () {
      const interfaceIdsToTest = [
        ...this.supportedInterfaces,
        DUMMY_UNSUPPORTED_ID,
      ];

      const supported = await this.mock.supportsAllInterfaces(
        this.target.address,
        interfaceIdsToTest,
      );
      expect(supported).to.equal(false);
    });
  });

  context('account address does not support KIP13', function () {
    it('does not support KIP13', async function () {
      const supported = await this.mock.supportsKIP13(DUMMY_ACCOUNT);
      expect(supported).to.equal(false);
    });

    it('does not support mock interface via supportsInterface', async function () {
      const supported = await this.mock.supportsInterface(
        DUMMY_ACCOUNT,
        DUMMY_ID,
      );
      expect(supported).to.equal(false);
    });

    it('does not support mock interface via supportsAllInterfaces', async function () {
      const supported = await this.mock.supportsAllInterfaces(
        DUMMY_ACCOUNT,
        [DUMMY_ID],
      );
      expect(supported).to.equal(false);
    });
  });
});
