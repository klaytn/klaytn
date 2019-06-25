// Copyright 2018 The klaytn Authors
// Copyright 2015 Jeffrey Wilcke, Felix Lange, Gustav Simonsson. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.
//
// This file is derived from crypto/secp256k1/ext.h (2018/06/04).
// See LICENSE in the top directory for the original copyright and license.

// secp256k1_context_create_sign_verify creates a context for signing and signature verification.
static secp256k1_context* secp256k1_context_create_sign_verify() {
	return secp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
}

// secp256k1_ext_ecdsa_recover recovers the public key of an encoded compact signature.
//
// Returns: 1: recovery was successful
//          0: recovery was not successful
// Args:    ctx:        pointer to a context object (cannot be NULL)
//  Out:    pubkey_out: the serialized 65-byte public key of the signer (cannot be NULL)
//  In:     sigdata:    pointer to a 65-byte signature with the recovery id at the end (cannot be NULL)
//          msgdata:    pointer to a 32-byte message (cannot be NULL)
static int secp256k1_ext_ecdsa_recover(
	const secp256k1_context* ctx,
	unsigned char *pubkey_out,
	const unsigned char *sigdata,
	const unsigned char *msgdata
) {
	secp256k1_ecdsa_recoverable_signature sig;
	secp256k1_pubkey pubkey;

	if (!secp256k1_ecdsa_recoverable_signature_parse_compact(ctx, &sig, sigdata, (int)sigdata[64])) {
		return 0;
	}
	if (!secp256k1_ecdsa_recover(ctx, &pubkey, &sig, msgdata)) {
		return 0;
	}
	size_t outputlen = 65;
	return secp256k1_ec_pubkey_serialize(ctx, pubkey_out, &outputlen, &pubkey, SECP256K1_EC_UNCOMPRESSED);
}

// secp256k1_ext_ecdsa_verify verifies an encoded compact signature.
//
// Returns: 1: signature is valid
//          0: signature is invalid
// Args:    ctx:        pointer to a context object (cannot be NULL)
//  In:     sigdata:    pointer to a 64-byte signature (cannot be NULL)
//          msgdata:    pointer to a 32-byte message (cannot be NULL)
//          pubkeydata: pointer to public key data (cannot be NULL)
//          pubkeylen:  length of pubkeydata
static int secp256k1_ext_ecdsa_verify(
	const secp256k1_context* ctx,
	const unsigned char *sigdata,
	const unsigned char *msgdata,
	const unsigned char *pubkeydata,
	size_t pubkeylen
) {
	secp256k1_ecdsa_signature sig;
	secp256k1_pubkey pubkey;

	if (!secp256k1_ecdsa_signature_parse_compact(ctx, &sig, sigdata)) {
		return 0;
	}
	if (!secp256k1_ec_pubkey_parse(ctx, &pubkey, pubkeydata, pubkeylen)) {
		return 0;
	}
	return secp256k1_ecdsa_verify(ctx, &sig, msgdata, &pubkey);
}

// secp256k1_ext_reencode_pubkey decodes then encodes a public key. It can be used to
// convert between public key formats. The input/output formats are chosen depending on the
// length of the input/output buffers.
//
// Returns: 1: conversion successful
//          0: conversion unsuccessful
// Args:    ctx:        pointer to a context object (cannot be NULL)
//  Out:    out:        output buffer that will contain the reencoded key (cannot be NULL)
//  In:     outlen:     length of out (33 for compressed keys, 65 for uncompressed keys)
//          pubkeydata: the input public key (cannot be NULL)
//          pubkeylen:  length of pubkeydata
static int secp256k1_ext_reencode_pubkey(
	const secp256k1_context* ctx,
	unsigned char *out,
	size_t outlen,
	const unsigned char *pubkeydata,
	size_t pubkeylen
) {
	secp256k1_pubkey pubkey;

	if (!secp256k1_ec_pubkey_parse(ctx, &pubkey, pubkeydata, pubkeylen)) {
		return 0;
	}
	unsigned int flag = (outlen == 33) ? SECP256K1_EC_COMPRESSED : SECP256K1_EC_UNCOMPRESSED;
	return secp256k1_ec_pubkey_serialize(ctx, out, &outlen, &pubkey, flag);
}

// secp256k1_ext_scalar_mul multiplies a point by a scalar in constant time.
//
// Returns: 1: multiplication was successful
//          0: scalar was invalid (zero or overflow)
// Args:    ctx:      pointer to a context object (cannot be NULL)
//  Out:    point:    the multiplied point (usually secret)
//  In:     point:    pointer to a 64-byte public point,
//                    encoded as two 256bit big-endian numbers.
//          scalar:   a 32-byte scalar with which to multiply the point
int secp256k1_ext_scalar_mul(const secp256k1_context* ctx, unsigned char *point, const unsigned char *scalar) {
	int ret = 0;
	int overflow = 0;
	secp256k1_fe feX, feY;
	secp256k1_gej res;
	secp256k1_ge ge;
	secp256k1_scalar s;
	ARG_CHECK(point != NULL);
	ARG_CHECK(scalar != NULL);
	(void)ctx;

	secp256k1_fe_set_b32(&feX, point);
	secp256k1_fe_set_b32(&feY, point+32);
	secp256k1_ge_set_xy(&ge, &feX, &feY);
	secp256k1_scalar_set_b32(&s, scalar, &overflow);
	if (overflow || secp256k1_scalar_is_zero(&s)) {
		ret = 0;
	} else {
		secp256k1_ecmult_const(&res, &ge, &s);
		secp256k1_ge_set_gej(&ge, &res);
		/* Note: can't use secp256k1_pubkey_save here because it is not constant time. */
		secp256k1_fe_normalize(&ge.x);
		secp256k1_fe_normalize(&ge.y);
		secp256k1_fe_get_b32(point, &ge.x);
		secp256k1_fe_get_b32(point+32, &ge.y);
		ret = 1;
	}
	secp256k1_scalar_clear(&s);
	return ret;
}

// eric.kim@groundx.xyz
// The following functions are added to support Schnorr signature scheme:
int secp256k1_ext_scalar_mul_bytes(const secp256k1_context* ctx, unsigned char *out, unsigned char *point, const unsigned char *scalar);
int secp256k1_ext_scalar_base_mult(const secp256k1_context* ctx, unsigned char *out, const unsigned char *scalar);
int secp256k1_ext_schnorr_verify(const secp256k1_context* ctx, unsigned char *P, unsigned char *R, const unsigned char *s, const unsigned char *e);
int secp256k1_ext_sc_mul(unsigned char *out, unsigned char *s1, unsigned char *s2);
int secp256k1_ext_sc_sub(unsigned char *out, unsigned char *s1, unsigned char *s2);
int secp256k1_ext_sc_add(unsigned char *out, unsigned char *s1, unsigned char *s2);

// secp256k1_ext_scalar_mul_bytes multiplies a point by a scalar in constant time.
//
// Returns: 1: multiplication was successful
//          0: scalar was invalid (zero or overflow)
// Args:    ctx:      pointer to a context object (cannot be NULL)
//  Out:    out:      the 64-byte multiplied point, encoded as two 256bit big-endian numbers
//  In:     point:    pointer to a 64-byte public point, encoded as two 256bit big-endian numbers
//          scalar:   a 32-byte scalar with which to multiply the point
int secp256k1_ext_scalar_mul_bytes(const secp256k1_context* ctx, unsigned char *out, unsigned char *point, const unsigned char *scalar) {
	int ret = 0;
	int overflow = 0;
	secp256k1_fe feX, feY;
	secp256k1_gej res;
	secp256k1_ge ge;
	secp256k1_scalar s;
	ARG_CHECK(point != NULL);
	ARG_CHECK(scalar != NULL);
	(void)ctx;

	secp256k1_fe_set_b32(&feX, point);
	secp256k1_fe_set_b32(&feY, point+32);
	secp256k1_ge_set_xy(&ge, &feX, &feY);
	secp256k1_scalar_set_b32(&s, scalar, &overflow);
	if (overflow || secp256k1_scalar_is_zero(&s)) {
		ret = 0;
	} else {
		secp256k1_ecmult_const(&res, &ge, &s);
		secp256k1_ge_set_gej(&ge, &res);
		/* Note: can't use secp256k1_pubkey_save here because it is not constant time. */
		secp256k1_fe_normalize(&ge.x);
		secp256k1_fe_normalize(&ge.y);
		secp256k1_fe_get_b32(out, &ge.x);
		secp256k1_fe_get_b32(out+32, &ge.y);
		ret = 1;
	}
	secp256k1_scalar_clear(&s);
	return ret;
}

// secp256k1_ext_scalar_base_mult multiplies a scalar to the SECP256k1 curve
// Out: res:    the 64-byte multiplied point, encoded as two 256-bit big-endian numbers
// In:  scalar: a 32-byte scalar with which to multiply the point
int secp256k1_ext_scalar_base_mult(const secp256k1_context* ctx, unsigned char *out, const unsigned char *scalar) {
    int overflow = 0;
    secp256k1_gej point;
    secp256k1_fe feX, feY;
    secp256k1_ge ge;
    secp256k1_scalar s;

    secp256k1_scalar_set_b32(&s, scalar, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&s)) {
        return 0;
    }

    secp256k1_ecmult_gen(&ctx->ecmult_gen_ctx, &point, &s);
    secp256k1_ge_set_gej(&ge, &point);
    /* Note: can't use secp256k1_pubkey_save here because it is not constant time. */
    secp256k1_fe_normalize(&ge.x);
    secp256k1_fe_normalize(&ge.y);
    secp256k1_fe_get_b32(out, &ge.x);
    secp256k1_fe_get_b32(out+32, &ge.y);
    secp256k1_scalar_clear(&s);
    return 1;
}

// secp256k1_ext_schnorr_verify verifies a Schnorr signature.
// Args:    ctx:    pointer to a context object (cannot be NULL)
//   In:    P:      a public key that is an elliptic curve point (64 bytes, big-endian)
//          R:      a part of a signature that is an elliptic curve point (64 bytes, big-endian)
//          s:      a part of a signature that is a scalar (32 bytes, big-endian)
//          e:      a derived random for this signature (32 bytes, big-endian; e = SHA256(msg || P || R))
int secp256k1_ext_schnorr_verify(const secp256k1_context* ctx, unsigned char *P, unsigned char *R,
                                    const unsigned char *s, const unsigned char *e) {
    int overflow = 0;
    secp256k1_fe feX, feY;
    secp256k1_gej V, ep, sg;
    secp256k1_ge ge;
    secp256k1_scalar sc;
    unsigned char tmp[64];
    (void)ctx;

    // compute e * P
    secp256k1_fe_set_b32(&feX, P);
    secp256k1_fe_set_b32(&feY, P+32);
    secp256k1_ge_set_xy(&ge, &feX, &feY);
    secp256k1_scalar_set_b32(&sc, e, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&sc)) {
        return 0;
    }
    secp256k1_ecmult_const(&ep, &ge, &sc); // ep => e * P
    secp256k1_scalar_clear(&sc);

    // compute s * G
    secp256k1_scalar_set_b32(&sc, s, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&sc)) {
        return 0;
    }
    secp256k1_ecmult_gen(&ctx->ecmult_gen_ctx, &sg, &sc);
    secp256k1_ge_set_gej(&ge, &sg); // ge => s * G
    secp256k1_scalar_clear(&sc);

    // compute s * G + e * P
    secp256k1_gej_add_ge(&V, &ep, &ge); // V => s * G + e * P

    secp256k1_ge_set_gej(&ge, &V);
    secp256k1_fe_normalize(&ge.x);
    secp256k1_fe_normalize(&ge.y);
    secp256k1_fe_get_b32(tmp, &ge.x);
    secp256k1_fe_get_b32(tmp+32, &ge.y); // tmp1 => s * G + e * P

    return 0 == memcmp(tmp, R, 64);
}

// secp256k1_ext_sc_mul multiplies two 32-byte scalars and returns the outcome.
// returns 0 in case of an overflow.
int secp256k1_ext_sc_mul(unsigned char *out, unsigned char *s1, unsigned char *s2) {
    int overflow = 0;
    secp256k1_scalar r, a, b;

    secp256k1_scalar_set_b32(&a, s1, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&a)) {
        return 0;
    }
    secp256k1_scalar_set_b32(&b, s2, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&b)) {
        return 0;
    }

    secp256k1_scalar_mul(&r, &a, &b);

    secp256k1_scalar_clear(&a);
    secp256k1_scalar_clear(&b);

    secp256k1_scalar_get_b32(out, &r);
    return 1;
}

// secp256k1_ext_sc_sub subtracts s2 from s1 where both s1 and s2 are 32-byte scalars.
// returns 0 in case of an overflow.
int secp256k1_ext_sc_sub(unsigned char *out, unsigned char *s1, unsigned char *s2) {
    int overflow = 0;
    secp256k1_scalar r, n, a, b;

    secp256k1_scalar_set_b32(&a, s1, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&a)) {
        return 0;
    }
    secp256k1_scalar_set_b32(&b, s2, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&b)) {
        return 0;
    }

    secp256k1_scalar_negate(&n, &b);
    secp256k1_scalar_add(&r, &a, &n);

    secp256k1_scalar_clear(&a);
    secp256k1_scalar_clear(&b);
    secp256k1_scalar_clear(&n);

    secp256k1_scalar_get_b32(out, &r);
    return 1;
}

// secp256k1_ext_sc_add adds s1 and s2 where both s1 and s2 are 32-byte scalars.
// returns 0 in case of an overflow.
int secp256k1_ext_sc_add(unsigned char *out, unsigned char *s1, unsigned char *s2) {
    int overflow = 0;
    secp256k1_scalar r, a, b;

    secp256k1_scalar_set_b32(&a, s1, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&a)) {
        return 0;
    }
    secp256k1_scalar_set_b32(&b, s2, &overflow);
    if (overflow || secp256k1_scalar_is_zero(&b)) {
        return 0;
    }

    secp256k1_scalar_add(&r, &a, &b);

    secp256k1_scalar_clear(&a);
    secp256k1_scalar_clear(&b);

    secp256k1_scalar_get_b32(out, &r);
    return 1;
}

